package readings

import (
	"database/sql"

	"github.com/ViBiOh/auth/auth"
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
  username = $1
`

func scanReadings(rows *sql.Rows) ([]*reading, error) {
	var (
		id     int64
		url    string
		public bool
		read   bool
	)

	list := make([]*reading, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &url, &public, &read); err != nil {
			return nil, err
		}

		list = append(list, &reading{id: id, URL: url, Public: public, Read: read})
	}

	return list, nil
}

func listReadingsOfUser(user *auth.User) ([]*reading, error) {
	rows, err := readingsDB.Query(listReadingsOfUserQuery, user.Username)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = db.RowsClose(`list readings`, rows, err)
	}()

	list, err := scanReadings(rows)
	if err != nil {
		return nil, err
	}

	return list, addTagsForReadings(list)
}
