package reading

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/eponae-api/pkg/tag"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
)

var _ crud.ItemService = &App{}

// App stores informations
type App struct {
	db                *sql.DB
	readingTagService *readingtag.App
	tagService        *tag.App
}

// NewService creates new ItemService
func NewService(db *sql.DB, readingTagService *readingtag.App, tagService *tag.App) *App {
	return &App{
		db:                db,
		readingTagService: readingTagService,
		tagService:        tagService,
	}
}

// Unmarsall a Reading
func (a App) Unmarsall(content []byte) (crud.Item, error) {
	var reading model.Reading

	if err := json.Unmarshal(content, &reading); err != nil {
		return nil, errors.WithStack(err)
	}

	return &reading, nil
}

// List readings of user
func (a App) List(ctx context.Context, page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]crud.Item, uint, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, 0, errors.New(`user not provided`)
	}

	list, total, err := a.listReadingsOfUser(user, page, pageSize, sortKey, sortAsc)
	if err != nil {
		return nil, 0, errors.Wrap(err, `unable to list readings of users`)
	}

	itemsList := make([]crud.Item, len(list))
	for index, item := range list {
		itemsList[index] = item
	}

	return itemsList, total, nil
}

// Get reading of user
func (a App) Get(ctx context.Context, ID string) (crud.Item, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New(`user not provided`)
	}

	reading, err := a.getReadingByID(user, ID)
	if err != nil {
		return nil, errors.Wrap(err, `unable to get reading`)
	}

	return reading, nil
}

// Create reading
func (a App) Create(ctx context.Context, o crud.Item) (item crud.Item, err error) {
	var reading *model.Reading
	reading, err = getReadingFromItem(ctx, o)
	if err != nil {
		return
	}

	if err = a.check(reading); err != nil {
		return nil, err
	}

	tx, err := a.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, `unable to get transaction`)
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	err = a.saveReading(reading, tx)
	if err != nil {
		err = errors.Wrap(err, `unable to create reading`)

		return
	}

	if err = a.readingTagService.CreateTagsForReading(reading, tx); err != nil {
		err = errors.Wrap(err, `unable to create reading's tags`)

		return
	}

	item = reading

	return
}

// Update reading
func (a App) Update(ctx context.Context, o crud.Item) (item crud.Item, err error) {
	var reading *model.Reading
	reading, err = getReadingFromItem(ctx, o)
	if err != nil {
		return
	}

	if err = a.check(reading); err != nil {
		return nil, err
	}

	tx, err := a.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, `unable to get transaction`)
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	err = a.saveReading(reading, tx)
	if err != nil {
		err = errors.Wrap(err, `unable to update reading`)

		return
	}

	if err = a.readingTagService.UpdateTagsForReading(reading, tx); err != nil {
		err = errors.Wrap(err, `unable to update reading's tags`)

		return
	}

	item = reading

	return
}

// Delete reading
func (a App) Delete(ctx context.Context, o crud.Item) (err error) {
	var reading *model.Reading
	reading, err = getReadingFromItem(ctx, o)
	if err != nil {
		return
	}

	err = a.deleteReading(reading, nil)
	if err != nil {
		err = errors.Wrap(err, `unable to delete reading`)
	}

	return
}

func (a App) check(o *model.Reading) error {
	if strings.TrimSpace(o.URL) == `` {
		return errors.Wrap(crud.ErrInvalid, `url is required`)
	}

	tagsIDs := make([]string, len(o.Tags))
	for index, tag := range o.Tags {
		tagsIDs[index] = tag.ID
	}

	tags, err := a.tagService.FindTagsByIds(tagsIDs)
	if err != nil {
		return errors.Wrap(err, `unable to list tags`)
	}

	if len(tags) != len(o.Tags) {
		return crud.ErrNotFound
	}

	return nil
}

func getReadingFromItem(ctx context.Context, o crud.Item) (*model.Reading, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New(`user not provided`)
	}

	item, ok := o.(*model.Reading)
	if !ok {
		return nil, errors.New(`item is not a reading`)
	}

	item.User = user

	return item, nil
}
