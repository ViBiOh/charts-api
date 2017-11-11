package readings

import (
	"database/sql"
	"fmt"

	"github.com/ViBiOh/httputils/db"
)

const listReadingsTagsOfReadingsQuery = `
SELECT
  readings_id,
  tags_id
FROM
  readings_tags
WHERE
  readings_id IN ($1)
ORDER BY
  readings_id ASC
`

type readingsTags struct {
	readingID uint
	tagID     uint
}

func scanReadingsTags(rows *sql.Rows, pageSize uint) ([]*readingsTags, error) {
	var (
		readingID uint
		tagID     uint
	)

	list := make([]*readingsTags, pageSize)

	for rows.Next() {
		if err := rows.Scan(&readingID, &tagID); err != nil {
			return nil, fmt.Errorf(`Error while scanning line: %v`, err)
		}

		list = append(list, &readingsTags{readingID, tagID})
	}

	return list, nil
}

func findReadingsTagsByIds(ids []uint) ([]*readingsTags, error) {
	rows, err := readingsDB.Query(listReadingsTagsOfReadingsQuery, db.WhereInUint(ids))
	if err != nil {
		return nil, fmt.Errorf(`Error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanReadingsTags(rows, uint(len(ids)))
}
