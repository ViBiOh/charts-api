package charts

import (
	"fmt"

	"github.com/ViBiOh/httputils/db"
)

const conservatoriesCountLabel = `conservatoriesCount`
const conservatoriesCountQuery = `
SELECT
  COUNT(id)
FROM
  conservatories
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
ORDER BY
  $1 %s
LIMIT
  $2
OFFSET
  $3
`

const conservatoriesSearchLabel = `conservatoriesSearch`
const conservatoriesSearchQuery = `
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
WHERE
  to_tsvector('french', name) @@ to_tsquery('french', $1)
  OR to_tsvector('french', city) @@ to_tsquery('french', $1)
  OR to_tsvector('french', zip) @@ to_tsquery('french', $1)
LIMIT
  $2
OFFSET
  $3
`

func countConservatories() (int64, error) {
	var count int64
	err := chartsDB.QueryRow(conservatoriesCountQuery).Scan(&count)

	return count, err
}

func readConservatories(page, pageSize int64, sortKey string, sortAsc bool) (conservatories []*conservatory, err error) {
	var offset int64
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	sortOrder := `ASC`
	if !sortAsc {
		sortOrder = `DESC`
	}

	rows, err := chartsDB.Query(fmt.Sprintf(conservatoriesQuery, sortOrder), sortKey, pageSize, offset)
	if err != nil {
		return
	}

	defer func() {
		err = db.RowsClose(conservatoriesLabel, rows, err)
	}()

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

	for rows.Next() {
		if err = rows.Scan(&id, &name, &category, &street, &city, &department, &zip, &latitude, &longitude); err != nil {
			return
		}

		conservatories = append(conservatories, &conservatory{ID: id, Name: name, Category: category, Street: street, City: city, Department: department, Zip: zip, Latitude: latitude, Longitude: longitude})
	}

	return
}
