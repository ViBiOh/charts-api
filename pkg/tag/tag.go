package tag

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ViBiOh/auth/v2/pkg/handler"
	authModel "github.com/ViBiOh/auth/v2/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
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
		return nil, err
	}

	return &tag, nil
}

// List tags of user
func (a App) List(ctx context.Context, page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]crud.Item, uint, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, 0, errors.New("user not provided")
	}

	list, total, err := a.listTagsOfUser(user, page, pageSize, sortKey, sortAsc)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to list tags of users: %w", err)
	}

	itemsList := make([]crud.Item, len(list))
	for index, item := range list {
		itemsList[index] = item
	}

	return itemsList, total, nil
}

// Get tag of user
func (a App) Get(ctx context.Context, ID uint64) (crud.Item, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, errors.New("user not provided")
	}

	tag, err := a.getTagByID(user, ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get tag: %w", err)
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
		err = fmt.Errorf("unable to create tag: %w", err)

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
		err = fmt.Errorf("unable to update tag: %w", err)

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
		err = fmt.Errorf("unable to delete tag: %w", err)
	}

	return
}

func (a App) check(o *model.Tag) error {
	if strings.TrimSpace(o.Name) == "" {
		return fmt.Errorf("name is required: %w", crud.ErrInvalid)
	}

	return nil
}

func getTagFromItem(ctx context.Context, o crud.Item) (*model.Tag, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, errors.New("user not provided")
	}

	item, ok := o.(*model.Tag)
	if !ok {
		return nil, errors.New("item is not a tag")
	}

	item.User = user

	return item, nil
}
