package charts

import "github.com/ViBiOh/httputils/db"

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
  lat,
  lng
FROM
  conservatories
LIMIT
  $1
OFFSET
  $2
ORDER BY
  $3 $4
`

// ReadConservatories retrieves conservatories
func readConservatories(page, pageSize int64, sortKey string, sortAsc bool) (conservatories []*conservatory, err error) {
	var offset int64
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	sortOrder := `ASC`
	if !sortAsc {
		sortOrder = `DESC`
	}

	rows, err := chartsDB.Query(conservatoriesLabel, conservatoriesQuery, pageSize, offset, sortKey, sortOrder)
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
