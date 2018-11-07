package tag

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/pkg/errors"
)

var _ crud.ItemService = &App{}

// App stores informations and secret of API
type App struct {
	db *sql.DB
}

// NewService creates new ItemService
func NewService(db *sql.DB) *App {
	return &App{
		db: db,
	}
}

// Unmarsall a Tag
func (a App) Unmarsall(content []byte) (crud.Item, error) {
	var tag model.Tag

	if err := json.Unmarshal(content, &tag); err != nil {
		return nil, errors.WithStack(err)
	}

	return &tag, nil
}

// List tags of user
func (a App) List(ctx context.Context, page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]crud.Item, uint, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, 0, errors.New(`user not provided`)
	}

	list, total, err := a.listTagsOfUser(user, page, pageSize, sortKey, sortAsc)
	if err != nil {
		logger.Error(`%+v`, err)
		return nil, 0, errors.New(`unable to list tags of users`)
	}

	itemsList := make([]crud.Item, len(list))
	for index, item := range list {
		itemsList[index] = item
	}

	return itemsList, total, nil
}

// Get tag of user
func (a App) Get(ctx context.Context, ID string) (crud.Item, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New(`user not provided`)
	}

	tag, err := a.getTagByID(user, ID)
	if err != nil {
		logger.Error(`%+v`, err)
		return nil, errors.New(`unable to get tag`)
	}

	return tag, nil
}

// Create tag
func (a App) Create(ctx context.Context, o crud.Item) (item crud.Item, err error) {
	var tag *model.Tag
	tag, err = getTagFromItem(ctx, o)
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

	tag.ID = ``

	err = a.saveTag(tag, tx)
	if err != nil {
		logger.Error(`%+v`, err)
		err = errors.New(`unable to create tag`)

		return
	}

	item = tag

	return
}

// Update tag
func (a App) Update(ctx context.Context, o crud.Item) (item crud.Item, err error) {
	var tag *model.Tag
	tag, err = getTagFromItem(ctx, o)
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

	err = a.saveTag(tag, nil)
	if err != nil {
		logger.Error(`%+v`, err)
		err = errors.New(`unable to update tag`)

		return
	}

	item = tag

	return
}

// Delete tag
func (a App) Delete(ctx context.Context, o crud.Item) (err error) {
	var tag *model.Tag
	tag, err = getTagFromItem(ctx, o)
	if err != nil {
		return
	}

	err = a.deleteTag(tag, nil)
	if err != nil {
		logger.Error(`%+v`, err)
		err = errors.New(`unable to delete tag`)
	}

	return
}

func getTagFromItem(ctx context.Context, o crud.Item) (*model.Tag, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New(`user not provided`)
	}

	item, ok := o.(*model.Tag)
	if !ok {
		return nil, errors.New(`item is not a tag`)
	}

	item.User = user

	return item, nil
}
