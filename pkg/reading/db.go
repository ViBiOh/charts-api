package reading

import (
	"database/sql"
	"fmt"

	authModel "github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/uuid"
)

func scanReading(row model.RowScanner) (*model.Reading, error) {
	var (
		id   string
		url  string
		read bool
	)

	err := row.Scan(&id, &url, &read)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, errors.WithStack(err)
	}

	return &model.Reading{UUID: id, URL: url, Read: read}, nil
}

func scanReadings(rows *sql.Rows) ([]*model.Reading, uint, error) {
	var (
		id         string
		url        string
		read       bool
		totalCount uint
	)

	list := make([]*model.Reading, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &url, &read, &totalCount); err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, err
			}

			return nil, 0, errors.WithStack(err)
		}

		list = append(list, &model.Reading{UUID: id, URL: url, Read: read})
	}

	return list, totalCount, nil
}

const listByUserQuery = `
SELECT
  id,
  url,
  read,
  count(*) OVER() AS full_count
FROM
  reading
WHERE
  user_id = $1
ORDER BY $4
LIMIT $2
OFFSET $3
`

func (a App) listReadingsOfUser(user *authModel.User, page, pageSize uint, sortKey string, sortAsc bool) ([]*model.Reading, uint, error) {
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

	list, totalCount, err := scanReadings(rows)
	if err != nil {
		return nil, 0, err
	}

	return list, totalCount, a.readingTagService.EnrichReadingsWithTags(list)
}

const getByIDQuery = `
SELECT
  id,
  user_id,
  url,
  read
FROM
  reading
WHERE
  id = $1
`

func (a App) getReadingByID(id string) (*model.Reading, error) {
	row := a.db.QueryRow(getByIDQuery, id)
	reading, err := scanReading(row)
	if err != nil {
		return nil, err
	}

	return reading, a.readingTagService.EnrichReadingWithTags(reading)
}

const insertQuery = `
INSERT INTO
  reading
(
  id,
  user_id,
  url,
  read
) VALUES (
  $1,
  $2,
  $3,
  $4
)
`

const updateQuery = `
UPDATE
  reading
SET
  url = $2,
  read = $3
WHERE
  id = $1
`

func (a App) saveReading(o *model.Reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New(`cannot save nil Reading`)
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

	if o.UUID != `` {
		if _, err = usedTx.Exec(updateQuery, o.UUID, o.URL, o.Read); err != nil {
			err = errors.WithStack(err)
		}
	} else {
		newID, err := uuid.New()
		if err != nil {
			return err
		}

		if _, err = usedTx.Exec(insertQuery, newID, o.User.ID, o.URL, o.Read); err != nil {
			err = errors.WithStack(err)
		} else {
			o.UUID = newID
		}
	}

	return
}

const deleteQuery = `
DELETE
  reading
WHERE
  id = $1
`

func (a App) deleteReading(o *model.Reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New(`cannot delete nil Reading`)
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

	if _, err = usedTx.Exec(deleteQuery, o.UUID); err != nil {
		err = errors.WithStack(err)
	}

	return
}
