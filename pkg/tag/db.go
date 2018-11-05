package tag

import (
	"database/sql"

	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/db"
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

	return &model.Tag{UUID: id, Name: name}, nil
}

func scanTags(rows *sql.Rows) ([]*model.Tag, error) {
	list := make([]*model.Tag, 0)

	for rows.Next() {
		Tag, err := scanTag(rows)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		list = append(list, Tag)
	}

	return list, nil
}

const listTagsByIDs = `
SELECT
  id,
  name
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

	return scanTags(rows)
}
