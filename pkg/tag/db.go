package tag

import (
	"database/sql"
	"fmt"

	authModel "github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func scanTag(row model.RowScanner) (*model.Tag, error) {
	var (
		id   string
		name string
	)

	err := row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, errors.WithStack(err)
	}

	return &model.Tag{ID: id, Name: name}, nil
}

func scanTags(rows *sql.Rows) ([]*model.Tag, uint, error) {
	var (
		id         string
		name       string
		totalCount uint
	)

	list := make([]*model.Tag, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &name, &totalCount); err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, err
			}

			return nil, 0, errors.WithStack(err)
		}

		list = append(list, &model.Tag{ID: id, Name: name})
	}

	return list, totalCount, nil
}

const listTagsByIDs = `
SELECT
  id,
  name,
  count(*) OVER() AS full_count
FROM
  tag
WHERE
  id = ANY ($1)
ORDER BY
  id
`

// FindTagsByIds finds tags by ids
func (a App) FindTagsByIds(ids []string) ([]*model.Tag, error) {
	rows, err := a.db.Query(listTagsByIDs, pq.Array(ids))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	tags, _, err := scanTags(rows)
	return tags, err
}

const listByUserQuery = `
SELECT
  id,
  name,
  count(*) OVER() AS full_count
FROM
  tag
WHERE
  user_id = $1
ORDER BY $4
LIMIT $2
OFFSET $3
`

func (a App) listTagsOfUser(user *authModel.User, page, pageSize uint, sortKey string, sortAsc bool) ([]*model.Tag, uint, error) {
	order := `creation_date DESC`

	if sortKey != `` {
		order = sortKey
	}
	if !sortAsc {
		order = fmt.Sprintf(`%s DESC`, order)
	}

	offset := (page - 1) * pageSize

	rows, err := a.db.Query(listByUserQuery, user.ID, pageSize, offset, order)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanTags(rows)
}

const getByIDQuery = `
SELECT
  id,
  name
FROM
  tag
WHERE
  user_id = $1
  AND id = $2
`

func (a App) getTagByID(user *authModel.User, id string) (*model.Tag, error) {
	return scanTag(a.db.QueryRow(getByIDQuery, user.ID, id))
}

const insertQuery = `
INSERT INTO
  tag
(
  id,
  user_id,
  name
) VALUES (
  $2,
  $1,
  $3
)
`

const updateQuery = `
UPDATE
  tag
SET
  name = $3
WHERE
  user_id = $1
  AND id = $2
`

func (a App) saveTag(o *model.Tag, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New(`cannot save nil Tag`)
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

	if o.ID != `` {
		if _, err = usedTx.Exec(updateQuery, o.User.ID, o.ID, o.Name); err != nil {
			err = errors.WithStack(err)
		}
	} else {
		var newID string
		newID, err = uuid.New()
		if err != nil {
			return err
		}

		if _, err = usedTx.Exec(insertQuery, o.User.ID, newID, o.Name); err != nil {
			err = errors.WithStack(err)
			return
		}

		o.ID = newID
	}

	return
}

const deleteQuery = `
DELETE FROM
  tag
WHERE
  user_id = $1
  AND id = $2
`

func (a App) deleteTag(o *model.Tag, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New(`cannot delete nil Tag`)
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

	if _, err = usedTx.Exec(deleteQuery, o.User.ID, o.ID); err != nil {
		err = errors.WithStack(err)
	}

	return
}
