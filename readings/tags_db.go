package readings

import (
	"database/sql"
	"fmt"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils/db"
)

const findTagsByidsQuery = `
SELECT
  id,
  name
FROM
  tags
WHERE
  id IN ($1)
`

const searchTagsWhereQuery = `
  AND to_tsvector('french', name) @@ to_tsquery('french', $INDEX)
`

const searchTagsCountQuery = `
SELECT
  COUNT(id)
FROM
  tags
WHERE
  user_id = $1
%s
`

const searchTagsQuery = `
SELECT
  id,
  name
FROM
  tags
WHERE
  user_id = $3
%s
ORDER BY
  %s %s
LIMIT
  $1
OFFSET
  $2
`

const readTagQuery = `
SELECT
  id,
  name
FROM
  tags
WHERE
  id = $1
  AND user_id = $2
`

const insertTagQuery = `
INSERT INTO
  tags
(
  user_id,
  name
) VALUES (
  $1,
  $2
)
RETURNING id
`

const updateTagQuery = `
UPDATE
  tags
SET
  name = $2
WHERE
  id = $1
`

const deleteTagQuery = `
DELETE FROM
  tags
WHERE
  id = $1
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

func scanTags(rows *sql.Rows, pageSize uint) ([]*tag, error) {
	var (
		id   uint
		name string
	)

	list := make([]*tag, 0, pageSize)

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf(`Error while scanning tag line: %v`, err)
		}

		list = append(list, &tag{ID: id, Name: name})
	}

	return list, nil
}

func searchTags(page, pageSize uint, sortKey string, sortAsc bool, user *auth.User, search string) ([]*tag, error) {
	if user == nil || user.ID == 0 {
		return nil, fmt.Errorf(`Unable to search tags of nil User`)
	}

	var offset uint
	if page > 1 {
		offset = (page - 1) * pageSize
	}

	sortOrder := `ASC`
	if !sortAsc {
		sortOrder = `DESC`
	}

	where, words := db.PrepareFullTextSearch(searchTagsWhereQuery, search, 4)

	var rows *sql.Rows
	var err error

	if where != `` {
		rows, err = readingsDB.Query(fmt.Sprintf(searchTagsQuery, where, sortKey, sortOrder), pageSize, offset, user.ID, words)
	} else {
		rows, err = readingsDB.Query(fmt.Sprintf(searchTagsQuery, where, sortKey, sortOrder), pageSize, offset, user.ID)
	}

	if err != nil {
		return nil, fmt.Errorf(`Error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanTags(rows, pageSize)
}

func countTags(user *auth.User, search string) (count uint, err error) {
	if user == nil || user.ID == 0 {
		return 0, fmt.Errorf(`Unable to count tags of nil User`)
	}

	where, words := db.PrepareFullTextSearch(searchTagsWhereQuery, search, 2)

	if where != `` {
		err = readingsDB.QueryRow(fmt.Sprintf(searchTagsCountQuery, where), user.ID, words).Scan(&count)
	} else {
		err = readingsDB.QueryRow(fmt.Sprintf(searchTagsCountQuery, where), user.ID).Scan(&count)
	}

	if err == sql.ErrNoRows {
		count = 0
		err = nil
	}

	if err != nil {
		err = fmt.Errorf(`Error while querying: %v`, err)
	}

	return
}

func findTagsByIds(ids []uint) ([]*tag, error) {
	rows, err := readingsDB.Query(findTagsByidsQuery, db.WhereInUint(ids))
	if err != nil {
		return nil, fmt.Errorf(`Error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanTags(rows, uint(len(ids)))
}

func getTag(id uint, user *auth.User) (*tag, error) {
	if user == nil {
		return nil, fmt.Errorf(`Unable to read tag of nil User`)
	}

	var (
		resultID uint
		name     string
	)

	if err := readingsDB.QueryRow(readTagQuery, id, user.ID).Scan(&resultID, &name); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf(`Error while querying: %v`, err)
	}

	return &tag{ID: resultID, Name: name, user: user}, nil
}

func saveTag(o *tag, tx *sql.Tx) (err error) {
	if o == nil {
		return fmt.Errorf(`Unable to save nil tag`)
	}

	if o.user == nil || o.user.ID == 0 {
		return fmt.Errorf(`Unable to save tag of nil User`)
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(readingsDB, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if o.ID != 0 {
		if _, err = usedTx.Exec(updateTagQuery, o.ID, o.Name); err != nil {
			err = fmt.Errorf(`Error while updating: %v`, err)
		}
	} else {
		var newID uint

		if err = usedTx.QueryRow(insertTagQuery, o.user.ID, o.Name).Scan(&newID); err != nil {
			err = fmt.Errorf(`Error while creating: %v`, err)
		} else {
			o.ID = newID
		}
	}

	return
}

func deleteTag(o *tag, tx *sql.Tx) (err error) {
	if o == nil || o.ID == 0 {
		return fmt.Errorf(`Unable to delete nil tag or one without ID`)
	}

	if o.user == nil || o.user.ID == 0 {
		return fmt.Errorf(`Unable to delete tag of nil User`)
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(readingsDB, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = usedTx.Exec(deleteTagQuery, o.ID); err != nil {
		err = fmt.Errorf(`Error while deleting: %v`, err)
	}

	return
}

func scanReadingsTagsForTag(rows *sql.Rows) (map[uint][]uint, error) {
	var (
		readingID uint
		tagID     uint
	)

	list := make(map[uint][]uint, 0)

	for rows.Next() {
		if err := rows.Scan(&readingID, &tagID); err != nil {
			return nil, fmt.Errorf(`Error while scanning line: %v`, err)
		}

		if _, ok := list[tagID]; ok {
			list[tagID] = append(list[tagID], readingID)
		} else {
			list[tagID] = []uint{readingID}
		}
	}

	return list, nil
}

func addTagsForReadings(readings []*reading) error {
	if len(readings) == 0 {
		return nil
	}

	ids := make([]uint, 0)
	for _, reading := range readings {
		ids = append(ids, reading.ID)
	}

	rows, err := readingsDB.Query(listReadingsTagsOfReadingsQuery, db.WhereInUint(ids))
	if err != nil {
		return fmt.Errorf(`Error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	tagLinks, err := scanReadingsTagsForTag(rows)
	if err != nil {
		return fmt.Errorf(`Error while scanning reading-tag of readings: %v`, err)
	} else if len(tagLinks) == 0 {
		return nil
	}

	tagsIds := make([]uint, 0)
	tagsByReading := make(map[uint][]uint, 0)
	for tagID, readingsIds := range tagLinks {
		tagsIds = append(tagsIds, tagID)

		for _, readingID := range readingsIds {
			if _, ok := tagsByReading[readingID]; ok {
				tagsByReading[readingID] = append(tagsByReading[readingID], tagID)
			} else {
				tagsByReading[readingID] = []uint{tagID}
			}
		}
	}

	tags, err := findTagsByIds(tagsIds)
	if err != nil {
		return fmt.Errorf(`Error while finding tags: %v`, err)
	}

	tagsByID := make(map[uint]*tag, 0)
	for _, tag := range tags {
		tagsByID[tag.ID] = tag
	}

	for _, reading := range readings {
		for _, tagID := range tagsByReading[reading.ID] {
			if reading.Tags == nil {
				reading.Tags = make([]*tag, 0)
			}

			reading.Tags = append(reading.Tags, tagsByID[tagID])
		}
	}

	return nil
}
