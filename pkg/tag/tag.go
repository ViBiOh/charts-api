package tag

import (
	"database/sql"

	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/httputils/pkg/crud"
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

// Empty creates an empty Tag
func (a App) Empty() crud.Item {
	return model.Tag{}
}

//List TODO
func (a *App) List(page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]crud.Item, error) {
	return nil, nil
}

//Get TODO
func (a *App) Get(ID string) (crud.Item, error) {
	return nil, nil
}

//Create TODO
func (a *App) Create(o crud.Item) (crud.Item, error) {
	return nil, nil
}

//Update TODO
func (a *App) Update(ID string, o crud.Item) (crud.Item, error) {
	return nil, nil
}

//Delete TODO
func (a *App) Delete(ID string) error {
	return nil
}
