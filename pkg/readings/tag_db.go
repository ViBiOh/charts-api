package readings

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/db"
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

var errNilUser = errors.New(`unable to request with nil User`)
var errNilTag = errors.New(`unable to save nil Tag`)

func scanTags(rows *sql.Rows, pageSize uint) ([]*tag, error) {
	var (
		id   uint
		name string
	)

	list := make([]*tag, 0, pageSize)

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf(`error while scanning tag line: %v`, err)
		}

		list = append(list, &tag{ID: id, Name: name})
	}

	return list, nil
}

func (a App) searchTags(page, pageSize uint, sortKey string, sortAsc bool, user *model.User, search string) ([]*tag, error) {
	if user == nil || user.ID == 0 {
		return nil, errNilUser
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
		rows, err = a.db.Query(fmt.Sprintf(searchTagsQuery, where, sortKey, sortOrder), pageSize, offset, user.ID, words)
	} else {
		rows, err = a.db.Query(fmt.Sprintf(searchTagsQuery, where, sortKey, sortOrder), pageSize, offset, user.ID)
	}

	if err != nil {
		return nil, fmt.Errorf(`error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanTags(rows, pageSize)
}

func (a App) countTags(user *model.User, search string) (count uint, err error) {
	if user == nil || user.ID == 0 {
		return 0, errNilUser
	}

	where, words := db.PrepareFullTextSearch(searchTagsWhereQuery, search, 2)

	if where != `` {
		err = a.db.QueryRow(fmt.Sprintf(searchTagsCountQuery, where), user.ID, words).Scan(&count)
	} else {
		err = a.db.QueryRow(fmt.Sprintf(searchTagsCountQuery, where), user.ID).Scan(&count)
	}

	if err == sql.ErrNoRows {
		count = 0
		err = nil
	}

	if err != nil {
		err = fmt.Errorf(`error while querying: %v`, err)
	}

	return
}

func (a App) findTagsByIds(ids []uint) ([]*tag, error) {
	rows, err := a.db.Query(findTagsByidsQuery, db.WhereInUint(ids))
	if err != nil {
		return nil, fmt.Errorf(`error while querying: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanTags(rows, uint(len(ids)))
}

func (a App) getTag(id uint, user *model.User) (*tag, error) {
	if user == nil {
		return nil, errNilUser
	}

	var (
		resultID uint
		name     string
	)

	if err := a.db.QueryRow(readTagQuery, id, user.ID).Scan(&resultID, &name); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf(`error while querying: %v`, err)
	}

	return &tag{ID: resultID, Name: name, user: user}, nil
}

func (a App) saveTag(o *tag, tx *sql.Tx) (err error) {
	if o == nil {
		return errNilTag
	}

	if o.user == nil || o.user.ID == 0 {
		return errNilUser
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(a.db, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if o.ID != 0 {
		if _, err = usedTx.Exec(updateTagQuery, o.ID, o.Name); err != nil {
			err = fmt.Errorf(`error while updating: %v`, err)
		}
	} else {
		var newID uint

		if err = usedTx.QueryRow(insertTagQuery, o.user.ID, o.Name).Scan(&newID); err != nil {
			err = fmt.Errorf(`error while creating: %v`, err)
		} else {
			o.ID = newID
		}
	}

	return
}

func (a App) deleteTag(o *tag, tx *sql.Tx) (err error) {
	if o == nil || o.ID == 0 {
		return errNilTag
	}

	if o.user == nil || o.user.ID == 0 {
		return errNilUser
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(a.db, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = usedTx.Exec(deleteTagQuery, o.ID); err != nil {
		err = fmt.Errorf(`error while deleting: %v`, err)
	}

	return
}
