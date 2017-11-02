package readings

import (
	"database/sql"
	"fmt"
)

const readUserQuery = `
SELECT
  id
FROM
  users
WHERE
  name = $1
`

func readUser(name string) (int64, error) {
	var (
		id int64
	)

	err := readingsDB.QueryRow(readUserQuery, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, fmt.Errorf(`Error while reading user: %v`, err)
	}

	return id, nil
}
