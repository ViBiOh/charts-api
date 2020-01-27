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

var _ crud.Service = &App{}

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

// Unmarshal a Tag
func (a App) Unmarshal(content []byte) (interface{}, error) {
	var tag model.Tag

	if err := json.Unmarshal(content, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// List tags of user
func (a App) List(ctx context.Context, page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]interface{}, uint, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, 0, model.ErrUserNotProvided
	}

	list, total, err := a.listTagsOfUser(user, page, pageSize, sortKey, sortAsc)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to list tags of users: %w", err)
	}

	itemsList := make([]interface{}, len(list))
	for index, item := range list {
		itemsList[index] = item
	}

	return itemsList, total, nil
}

// Get tag of user
func (a App) Get(ctx context.Context, ID uint64) (interface{}, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, model.ErrUserNotProvided
	}

	tag, err := a.getTagByID(user, ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get tag: %w", err)
	}

	return tag, nil
}

// Create tag
func (a App) Create(ctx context.Context, o interface{}) (item interface{}, err error) {
	var tag *model.Tag
	tag, err = getTagFromItem(ctx, o)
	if err != nil {
		return
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
func (a App) Update(ctx context.Context, o interface{}) (item interface{}, err error) {
	var tag *model.Tag
	tag, err = getTagFromItem(ctx, o)
	if err != nil {
		return
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
func (a App) Delete(ctx context.Context, o interface{}) (err error) {
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

// Check instance
func (a App) Check(_ context.Context, _, new interface{}) []crud.Error {
	if new == nil {
		return nil
	}

	tag := new.(*model.Tag)
	errors := make([]crud.Error, 0)

	if strings.TrimSpace(tag.Name) == "" {
		errors = append(errors, crud.NewError("name", "name is required"))
	}

	return errors
}

func getTagFromItem(ctx context.Context, o interface{}) (*model.Tag, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, model.ErrUserNotProvided
	}

	item, ok := o.(*model.Tag)
	if !ok {
		return nil, errors.New("item is not a tag")
	}

	item.User = user

	return item, nil
}
