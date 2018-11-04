package readings

import (
	"database/sql"

	"github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
)

const listReadingsOfUserQuery = `
SELECT
  id,
  url,
  read
FROM
  readings
WHERE
  user_id = $1
`

const insertReading = `
INSERT INTO
  readings
(
  url,
  user_id,
  read
) VALUES (
  $1,
  $2,
  $3
)
RETURNING id
`
const updateReading = `
UPDATE
  tags
SET
  url = $2,
  read = $3
WHERE
  id = $1
`

var errNilReading = errors.New(`unable to save nil reading`)

func scanReadings(rows *sql.Rows) ([]*reading, error) {
	var (
		id   uint
		url  string
		read bool
	)

	list := make([]*reading, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &url, &read); err != nil {
			return nil, errors.WithStack(err)
		}

		list = append(list, &reading{ID: id, URL: url, Read: read})
	}

	return list, nil
}

func (a App) listReadingsOfUser(user *model.User) ([]*reading, error) {
	rows, err := a.db.Query(listReadingsOfUserQuery, user.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	list, err := scanReadings(rows)
	if err != nil {
		return nil, err
	}

	return a.enrichReadingsWithTags(list)
}

func (a App) saveReading(o *reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errNilReading
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
		if _, err = usedTx.Exec(updateReading, o.ID, o.URL, o.Read); err != nil {
			err = errors.WithStack(err)
		}
	} else {
		var newID uint

		if err = usedTx.QueryRow(insertReading, o.user.ID, o.URL, o.Read).Scan(&newID); err != nil {
			err = errors.WithStack(err)
		} else {
			o.ID = newID
		}
	}

	return
}
