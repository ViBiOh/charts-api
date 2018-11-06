package tag

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/crud"
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

//List TODO
func (a App) List(ctx context.Context, page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]crud.Item, uint, error) {
	return nil, 0, nil
}

//Get TODO
func (a App) Get(ctx context.Context, ID string) (crud.Item, error) {
	return nil, nil
}

//Create TODO
func (a App) Create(ctx context.Context, o crud.Item) (crud.Item, error) {
	return nil, nil
}

//Update TODO
func (a App) Update(ctx context.Context, o crud.Item) (crud.Item, error) {
	return nil, nil
}

//Delete TODO
func (a App) Delete(ctx context.Context, o crud.Item) error {
	return nil
}
