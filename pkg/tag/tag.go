package tag

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/errors"
)

var _ crud.ItemService = &App{}

// App of package
type App struct {
	db *sql.DB
}

// New creates new App
func New(db *sql.DB) *App {
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
		return nil, 0, errors.Wrap(err, `unable to list tags of users`)
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
		return nil, errors.Wrap(err, `unable to get tag`)
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

	if err = a.check(tag); err != nil {
		return nil, err
	}

	err = a.saveTag(tag, nil)
	if err != nil {
		err = errors.Wrap(err, `unable to create tag`)

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

	if err = a.check(tag); err != nil {
		return nil, err
	}

	err = a.saveTag(tag, nil)
	if err != nil {
		err = errors.Wrap(err, `unable to update tag`)

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
		err = errors.Wrap(err, `unable to delete tag`)
	}

	return
}

func (a App) check(o *model.Tag) error {
	if strings.TrimSpace(o.Name) == `` {
		return errors.Wrap(crud.ErrInvalid, `name is required`)
	}

	return nil
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
