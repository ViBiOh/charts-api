package readings

import (
	"database/sql"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils/db"
)

const listTagsOfUserQuery = `
SELECT
  id,
  name
FROM
  tags
WHERE
  username = $1
`

const listTagsByidsQuery = `
SELECT
  id,
  name
FROM
  tags
WHERE
  id IN ($1)
`

const listReadingsTagsOfReadingsQuery = `
SELECT
  readings_id,
  tags_id
FROM
  readings_tags
WHERE
  readings_id IN ($1)
`

func scanTags(rows *sql.Rows) ([]*tag, error) {
	var (
		id   int64
		name string
	)

	list := make([]*tag, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		list = append(list, &tag{id: id, Name: name})
	}

	return list, nil
}

func scanReadingsTagsForTag(rows *sql.Rows) (map[int64][]int64, error) {
	var (
		readingID int64
		tagID     int64
	)

	list := make(map[int64][]int64, 0)

	for rows.Next() {
		if err := rows.Scan(&readingID, &tagID); err != nil {
			return nil, err
		}

		if _, ok := list[tagID]; ok {
			list[tagID] = append(list[tagID], readingID)
		} else {
			list[tagID] = []int64{readingID}
		}
	}

	return list, nil
}

func listTagsOfUser(user *auth.User) ([]*tag, error) {
	rows, err := readingsDB.Query(listTagsOfUserQuery, user.Username)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = db.RowsClose(`list readings`, rows, err)
	}()

	return scanTags(rows)
}

func listTagsByIds(ids []int64) ([]*tag, error) {
	rows, err := readingsDB.Query(listTagsByidsQuery, ids)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = db.RowsClose(`list tags by ids`, rows, err)
	}()

	return scanTags(rows)
}

func addTagsForReadings(readings []*reading) error {
	if len(readings) == 0 {
		return nil
	}

	ids := make([]int64, 0)
	for _, reading := range readings {
		ids = append(ids, reading.id)
	}

	rows, err := readingsDB.Query(listReadingsTagsOfReadingsQuery, ids)
	if err != nil {
		return err
	}

	defer func() {
		err = db.RowsClose(`list readings-tags of readings`, rows, err)
	}()

	tagLinks, err := scanReadingsTagsForTag(rows)
	if err != nil {
		return err
	} else if len(tagLinks) == 0 {
		return nil
	}

	tagsIds := make([]int64, 0)
	tagsByReading := make(map[int64][]int64, 0)
	for tagID, readingsIds := range tagLinks {
		tagsIds = append(tagsIds, tagID)

		for _, readingID := range readingsIds {
			if _, ok := tagsByReading[readingID]; ok {
				tagsByReading[readingID] = append(tagsByReading[readingID], tagID)
			} else {
				tagsByReading[readingID] = []int64{tagID}
			}
		}
	}

	tags, err := listTagsByIds(tagsIds)
	if err != nil {
		return err
	}

	tagsByID := make(map[int64]*tag, 0)
	for _, tag := range tags {
		tagsByID[tag.id] = tag
	}

	for _, reading := range readings {
		for _, tagID := range tagsByReading[reading.id] {
			if reading.Tags == nil {
				reading.Tags = make([]*tag, 0)
			}

			reading.Tags = append(reading.Tags, tagsByID[tagID])
		}
	}

	return nil
}
