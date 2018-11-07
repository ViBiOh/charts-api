package reading

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/logger"
)

var _ crud.ItemService = &App{}

// App stores informations
type App struct {
	db                *sql.DB
	readingTagService *readingtag.App
}

// NewService creates new ItemService
func NewService(db *sql.DB, readingTagService *readingtag.App) *App {
	return &App{
		db:                db,
		readingTagService: readingTagService,
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
		logger.Error(`%+v`, err)
		return nil, 0, errors.New(`unable to list readings of users`)
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
		logger.Error(`%+v`, err)
		return nil, errors.New(`unable to get reading`)
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

	tx, err := a.db.Begin()
	if err != nil {
		logger.Error(`%+v`, err)
		return nil, errors.New(`unable to get transaction`)
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	reading.ID = ``

	err = a.saveReading(reading, tx)
	if err != nil {
		logger.Error(`%+v`, err)
		err = errors.New(`unable to create reading`)

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

	tx, err := a.db.Begin()
	if err != nil {
		logger.Error(`%+v`, err)
		return nil, errors.New(`unable to get transaction`)
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	err = a.saveReading(reading, nil)
	if err != nil {
		logger.Error(`%+v`, err)
		err = errors.New(`unable to update reading`)

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
		logger.Error(`%+v`, err)
		err = errors.New(`unable to delete reading`)
	}

	return
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
