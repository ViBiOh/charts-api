package reading

import (
	"database/sql"

	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/httputils/pkg/crud"
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

// Empty creates an empty Reading
func (a *App) Empty() crud.Item {
	return model.Reading{}
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
