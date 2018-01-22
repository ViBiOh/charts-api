package readings

import (
	"database/sql"
	"errors"
	"fmt"

	authProvider "github.com/ViBiOh/auth/provider"
	"github.com/ViBiOh/httputils/db"
)

const listReadingsOfUserQuery = `
SELECT
  id,
  url,
  public,
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
  public,
  read
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING id
`
const updateReading = `
UPDATE
  tags
SET
  url = $2,
  public = $3,
  read = $4
WHERE
  id = $1
`

var errNilReading = errors.New(`Unable to save nil reading`)

func scanReadings(rows *sql.Rows) ([]*reading, error) {
	var (
		id     uint
		url    string
		public bool
		read   bool
	)

	list := make([]*reading, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &url, &public, &read); err != nil {
			return nil, fmt.Errorf(`Error while scanning reading line: %v`, err)
		}

		list = append(list, &reading{ID: id, URL: url, Public: public, Read: read})
	}

	return list, nil
}

func listReadingsOfUser(user *authProvider.User) ([]*reading, error) {
	rows, err := readingsDB.Query(listReadingsOfUserQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf(`Error while listing readings of user: %v`, err)
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	list, err := scanReadings(rows)
	if err != nil {
		return nil, fmt.Errorf(`Error while scanning readings: %v`, err)
	}

	return enrichReadingsWithTags(list)
}

func saveReading(o *reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errNilReading
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
		if _, err = usedTx.Exec(updateReading, o.ID, o.URL, o.Public, o.Read); err != nil {
			err = fmt.Errorf(`Error while updating reading for user=%s: %v`, o.user.Username, err)
		}
	} else {
		var newID uint

		if err = usedTx.QueryRow(insertReading, o.user.ID, o.URL, o.Public, o.Read).Scan(&newID); err != nil {
			err = fmt.Errorf(`Error while creating reading for user=%s: %v`, o.user.Username, err)
		} else {
			o.ID = newID
		}
	}

	return
}
