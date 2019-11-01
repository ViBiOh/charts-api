package reading

import (
	"database/sql"
	"errors"
	"fmt"

	authModel "github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/uuid"
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

		return nil, err
	}

	return &model.Reading{ID: id, URL: url, Read: read}, nil
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

			return nil, 0, err
		}

		list = append(list, &model.Reading{ID: id, URL: url, Read: read})
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
	order := "creation_date DESC"

	if sortKey != "" {
		order = sortKey
	}
	if !sortAsc {
		order = fmt.Sprintf("%s DESC", order)
	}

	offset := (page - 1) * pageSize

	rows, err := a.db.Query(listByUserQuery, user.ID, pageSize, offset, order)
	if err != nil {
		return nil, 0, err
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
  url,
  read
FROM
  reading
WHERE
  user_id = $1
  AND id = $2
`

func (a App) getReadingByID(user *authModel.User, id string) (*model.Reading, error) {
	row := a.db.QueryRow(getByIDQuery, user.ID, id)
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
  $2,
  $1,
  $3,
  $4
)
`

const updateQuery = `
UPDATE
  reading
SET
  url = $3,
  read = $4
WHERE
  user_id = $1
  AND id = $2
`

func (a App) saveReading(o *model.Reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New("cannot save nil Reading")
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

	if o.ID != "" {
		_, err = usedTx.Exec(updateQuery, o.User.ID, o.ID, o.URL, o.Read)
	} else {
		var newID string
		newID, err = uuid.New()
		if err != nil {
			return err
		}

		_, err = usedTx.Exec(insertQuery, o.User.ID, newID, o.URL, o.Read)
		if err != nil {
			return
		}

		o.ID = newID
	}

	if err = a.readingTagService.SaveTagsForReading(o, usedTx); err != nil {
		err = fmt.Errorf("unable to save reading's tags: %w", err)

		return
	}

	return
}

const deleteQuery = `
DELETE FROM
  reading
WHERE
  user_id = $1
  AND id = $2
`

func (a App) deleteReading(o *model.Reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New("cannot delete nil Reading")
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

	_, err = usedTx.Exec(deleteQuery, o.User.ID, o.ID)
	return
}
