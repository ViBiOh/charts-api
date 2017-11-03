package readings

import (
	"database/sql"
	"fmt"

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
  user_id = $1
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
			return nil, fmt.Errorf(`Error while scanning reading line: %v`, err)
		}

		list = append(list, &reading{id: id, URL: url, Public: public, Read: read})
	}

	return list, nil
}

func listReadingsOfUser(user *auth.User) ([]*reading, error) {
	rows, err := readingsDB.Query(listReadingsOfUserQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf(`Error while listing readings of user: %v`, err)
	}

	defer func() {
		err = db.RowsClose(`list readings of user`, rows, err)
	}()

	list, err := scanReadings(rows)
	if err != nil {
		return nil, fmt.Errorf(`Error while scanning readings: %v`, err)
	}

	return list, addTagsForReadings(list)
}
