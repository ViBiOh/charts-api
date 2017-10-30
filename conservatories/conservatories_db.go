package conservatories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ViBiOh/httputils/db"
)

const conservatoriesCountLabel = `conservatoriesCount`
const conservatoriesCountQuery = `
SELECT
  COUNT(id)
FROM
  conservatories
%s
`

const conservatoriesLabel = `conservatories`
const conservatoriesQuery = `
SELECT
  id,
  name,
  category,
  street,
  city,
  department,
  zip,
  latitude,
  longitude
FROM
	conservatories
%s
ORDER BY
  %s %s
LIMIT
  $1
OFFSET
  $2
`

const conservatoriesSearchWhere = `
WHERE
  to_tsvector('french', name) @@ to_tsquery('french', $INDEX)
  OR to_tsvector('french', city) @@ to_tsquery('french', $INDEX)
  OR to_tsvector('french', zip) @@ to_tsquery('french', $INDEX)
`

const conservatoriesByDepartementLabel = `conservatoriesByDepartement`
const conservatoriesByDepartementQuery = `
SELECT
  department,
  COUNT(id)
FROM
  conservatories
GROUP BY
  department
`

const conservatoriesByZipOfDepartmentQuery = `
SELECT
  zip,
  COUNT(id)
FROM
  conservatories
WHERE
  department = $1
GROUP BY
  zip
`

func scanConservatoryRows(rows *sql.Rows, pageSize int64) ([]*conservatory, error) {
	var (
		id         int64
		name       string
		category   string
		street     string
		city       string
		department int
		zip        string
		latitude   float64
		longitude  float64
	)

	conservatories := make([]*conservatory, 0, pageSize)

	for rows.Next() {
		if err := rows.Scan(&id, &name, &category, &street, &city, &department, &zip, &latitude, &longitude); err != nil {
			return nil, err
		}

		conservatories = append(conservatories, &conservatory{ID: id, Name: name, Category: category, Street: street, City: city, Department: department, Zip: zip, Latitude: latitude, Longitude: longitude})
	}

	return conservatories, nil
}

func scanAggregateRows(rows *sql.Rows) (map[string]int64, error) {
	var (
		key   string
		count int64
	)

	aggregate := make(map[string]int64, 0)

	for rows.Next() {
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}

		aggregate[key] = count
	}

	return aggregate, nil
}

func prepareFullTextSearch(search string, index int) (string, string) {
	if search == `` {
		return ``, ``
	}

	words := strings.Split(search, ` `)
	transformedWords := make([]string, 0, len(words))

	for _, word := range words {
		transformedWords = append(transformedWords, word+`:*`)
	}

	return strings.Replace(conservatoriesSearchWhere, `$INDEX`, fmt.Sprintf(`$%d`, index), -1), strings.Join(transformedWords, ` | `)
}

func countConservatories(search string) (count int64, err error) {
	where, words := prepareFullTextSearch(search, 1)

	if words != `` {
		err = chartsDB.QueryRow(fmt.Sprintf(conservatoriesCountQuery, where), words).Scan(&count)
	} else {
		err = chartsDB.QueryRow(fmt.Sprintf(conservatoriesCountQuery, where)).Scan(&count)
	}

	return
}

func searchConservatories(page, pageSize int64, sortKey string, sortAsc bool, search string) ([]*conservatory, error) {
	var offset int64
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	sortOrder := `ASC`
	if !sortAsc {
		sortOrder = `DESC`
	}

	where, words := prepareFullTextSearch(search, 3)
	var rows *sql.Rows
	var err error

	if words != `` {
		rows, err = chartsDB.Query(fmt.Sprintf(conservatoriesQuery, where, sortKey, sortOrder), pageSize, offset, words)
	} else {
		rows, err = chartsDB.Query(fmt.Sprintf(conservatoriesQuery, where, sortKey, sortOrder), pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	defer func() {
		err = db.RowsClose(conservatoriesLabel, rows, err)
	}()

	return scanConservatoryRows(rows, pageSize)
}

func countByDepartment() (map[string]int64, error) {
	rows, err := chartsDB.Query(conservatoriesByDepartementQuery)

	if err != nil {
		return nil, err
	}
	return scanAggregateRows(rows)
}

func countByZipOfDepartment(department string) (map[string]int64, error) {
	rows, err := chartsDB.Query(conservatoriesByZipOfDepartmentQuery, department)

	if err != nil {
		return nil, err
	}
	return scanAggregateRows(rows)
}
