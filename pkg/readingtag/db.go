package readingtag

import (
	"database/sql"
	"errors"

	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/lib/pq"
)

func scanReadingTag(row model.RowScanner) (*model.ReadingTag, error) {
	var (
		readingID uint64
		tagID     uint64
	)

	err := row.Scan(&readingID, &tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	return &model.ReadingTag{ReadingID: readingID, TagID: tagID}, nil
}

func scanReadingTags(rows *sql.Rows) ([]*model.ReadingTag, error) {
	list := make([]*model.ReadingTag, 0)

	for rows.Next() {
		readingTag, err := scanReadingTag(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, readingTag)
	}

	return list, nil
}

const listTagsByReadingIDsQuery = `
SELECT
  reading_id,
  tag_id
FROM
  reading_tag
WHERE
  reading_id = ANY ($1)
ORDER BY
  reading_id ASC
`

func (a App) listTagsByReadingIDs(ids []uint64) ([]*model.ReadingTag, error) {
	rows, err := a.db.Query(listTagsByReadingIDsQuery, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer func() {
		err = db.RowsClose(rows, err)
	}()

	return scanReadingTags(rows)
}

// EnrichReadingsWithTags update given reading with tags data
func (a App) EnrichReadingsWithTags(readings []*model.Reading) error {
	if len(readings) == 0 {
		return nil
	}

	readingsID := make([]uint64, len(readings))
	for i, reading := range readings {
		readingsID[i] = reading.ID
	}

	readingTags, err := a.listTagsByReadingIDs(readingsID)
	if err != nil {
		return err
	}

	tagsID := make([]uint64, len(readingTags))
	tagsByReading := make(map[uint64][]uint64, 0)
	for i, link := range readingTags {
		tagsID[i] = link.TagID

		if e, ok := tagsByReading[link.ReadingID]; ok {
			tagsByReading[link.ReadingID] = append(e, link.TagID)
		} else {
			tagsByReading[link.ReadingID] = []uint64{link.TagID}
		}
	}

	tags, err := a.tagService.FindTagsByIds(tagsID)
	if err != nil {
		return err
	}

	tagsByID := make(map[uint64]*model.Tag, 0)
	for _, tagObj := range tags {
		if _, ok := tagsByID[tagObj.ID]; !ok {
			tagsByID[tagObj.ID] = tagObj
		}
	}

	for _, reading := range readings {
		if tagsID, ok := tagsByReading[reading.ID]; ok {
			for _, tagID := range tagsID {
				if tagObj, ok := tagsByID[tagID]; ok {
					if reading.Tags == nil {
						reading.Tags = []*model.Tag{tagObj}
					} else {
						reading.Tags = append(reading.Tags, tagObj)
					}
				}
			}
		}
	}

	return nil
}

// EnrichReadingWithTags update given reading with tags data
func (a App) EnrichReadingWithTags(o *model.Reading) error {
	return a.EnrichReadingsWithTags([]*model.Reading{o})
}

// SaveTagsForReading update tags' link for given reading
func (a App) SaveTagsForReading(o *model.Reading, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New("cannot create tag for nil Reading")
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

	newTagIDs := make([]uint64, len(o.Tags))
	for index, tag := range o.Tags {
		newTagIDs[index] = tag.ID
	}

	var existingTags []*model.ReadingTag
	existingTags, err = a.listTagsByReadingIDs([]uint64{o.ID})
	if err != nil {
		return
	}

	existingTagIDs := make([]uint64, len(existingTags))
	for index, existingTag := range existingTags {
		existingTagIDs[index] = existingTag.TagID

		if !IncludesUint64(newTagIDs, existingTag.TagID) {
			err = a.deleteReadingTag(existingTag, usedTx)
			if err != nil {
				return
			}
		}
	}

	for _, newTagID := range newTagIDs {
		if !IncludesUint64(existingTagIDs, newTagID) {
			if err = a.insertReadingTag(&model.ReadingTag{ReadingID: o.ID, TagID: newTagID}, usedTx); err != nil {
				return
			}
		}
	}

	return
}

const insertQuery = `
INSERT INTO
  reading_tag
(
  reading_id,
  tag_id
) VALUES (
  $1,
  $2
)
`

func (a App) insertReadingTag(o *model.ReadingTag, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New("cannot insert nil ReadingTag")
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

	_, err = usedTx.Exec(insertQuery, o.ReadingID, o.TagID)
	return
}

const deleteQuery = `
DELETE FROM
  reading_tag
WHERE
  reading_id = $1
  AND tag_id = $2
`

func (a App) deleteReadingTag(o *model.ReadingTag, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New("cannot delete nil ReadingTag")
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

	_, err = usedTx.Exec(deleteQuery, o.ReadingID, o.TagID)
	return
}
