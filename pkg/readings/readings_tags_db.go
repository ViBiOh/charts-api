package readings

import (
	"database/sql"
	"fmt"

	"github.com/ViBiOh/httputils/pkg/db"
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

func (a *App) scanReadingsTags(rows *sql.Rows, pageSize uint) ([]*readingsTags, error) {
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

func (a *App) findReadingsTagsByIds(ids []uint) ([]*readingsTags, error) {
	rows, err := a.db.Query(listReadingsTagsOfReadingsQuery, db.WhereInUint(ids))
	if err != nil {
		return nil, fmt.Errorf(`Error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return a.scanReadingsTags(rows, uint(len(ids)))
}

func (a *App) enrichReadingsWithTags(readings []*reading) ([]*reading, error) {
	if len(readings) == 0 {
		return readings, nil
	}

	readingsID := make([]uint, len(readings))
	for i, reading := range readings {
		readingsID[i] = reading.ID
	}

	tagsLink, err := a.findReadingsTagsByIds(readingsID)
	if err != nil {
		return nil, fmt.Errorf(`Error while finding readings tags: %v`, err)
	}

	tagsID := make([]uint, len(tagsLink))
	tagsByReading := make(map[uint][]uint, 0)
	for i, link := range tagsLink {
		tagsID[i] = link.tagID

		if e, ok := tagsByReading[link.readingID]; ok {
			tagsByReading[link.readingID] = append(e, link.tagID)
		} else {
			tagsByReading[link.readingID] = []uint{link.tagID}
		}
	}

	tags, err := a.findTagsByIds(tagsID)
	if err != nil {
		return nil, fmt.Errorf(`Error while finding tags: %v`, err)
	}

	tagsByID := make(map[uint]*tag, 0)
	for _, tagObj := range tags {
		if _, ok := tagsByID[tagObj.ID]; !ok {
			tagsByID[tagObj.ID] = tagObj
		}
	}

	for _, reading := range readings {
		if tagsID, ok := tagsByReading[reading.ID]; ok {
			for _, tagID := range tagsID {
				if tagObj, ok := tagsByID[tagID]; ok {
					if reading.Tags == nil {
						reading.Tags = []*tag{tagObj}
					} else {
						reading.Tags = append(reading.Tags, tagObj)
					}
				}
			}
		}
	}

	return readings, nil
}
