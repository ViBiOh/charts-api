package reading

import (
	"context"
	"database/sql"

	"github.com/ViBiOh/auth/pkg/auth"
	"github.com/ViBiOh/eponae-api/pkg/model"
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/httputils/pkg/crud"
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

// Empty creates an empty Reading
func (a App) Empty() crud.Item {
	return model.Reading{}
}

//List TODO
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
