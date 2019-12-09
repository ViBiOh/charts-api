package reading

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
	"github.com/ViBiOh/eponae-api/pkg/readingtag"
	"github.com/ViBiOh/eponae-api/pkg/tag"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
)

var _ crud.Service = &App{}

// App of package
type App struct {
	db                *sql.DB
	readingTagService *readingtag.App
	tagService        *tag.App
}

// New creates new App
func New(db *sql.DB, readingTagService *readingtag.App, tagService *tag.App) *App {
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
		return nil, err
	}

	return &reading, nil
}

// List readings of user
func (a App) List(ctx context.Context, page, pageSize uint, sortKey string, sortAsc bool, filters map[string][]string) ([]crud.Item, uint, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, 0, errors.New("user not provided")
	}

	list, total, err := a.listReadingsOfUser(user, page, pageSize, sortKey, sortAsc)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to list readings of users: %w", err)
	}

	itemsList := make([]crud.Item, len(list))
	for index, item := range list {
		itemsList[index] = item
	}

	return itemsList, total, nil
}

// Get reading of user
func (a App) Get(ctx context.Context, ID uint64) (crud.Item, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, errors.New("user not provided")
	}

	reading, err := a.getReadingByID(user, ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get reading: %w", err)
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

	err = a.saveReading(reading, nil)
	if err != nil {
		err = fmt.Errorf("unable to create reading: %w", err)

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

	err = a.saveReading(reading, nil)
	if err != nil {
		err = fmt.Errorf("unable to update reading: %w", err)

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
		err = fmt.Errorf("unable to delete reading: %w", err)
	}

	return
}

// Check instance
func (a App) Check(o crud.Item) []error {
	reading := o.(*model.Reading)
	errors := make([]error, 0)

	if strings.TrimSpace(reading.URL) == "" {
		errors = append(errors, fmt.Errorf("url is required: %w", crud.ErrInvalid))
	}

	tagsIDs := make([]uint64, len(reading.Tags))
	for index, tag := range reading.Tags {
		tagsIDs[index] = tag.ID
	}

	tags, err := a.tagService.FindTagsByIds(tagsIDs)
	if err != nil {
		errors = append(errors, fmt.Errorf("unable to list tags: %w", err))
	}

	if len(tags) != len(reading.Tags) {
		errors = append(errors, crud.ErrNotFound)
	}

	return errors
}

func getReadingFromItem(ctx context.Context, o crud.Item) (*model.Reading, error) {
	user := handler.UserFromContext(ctx)
	if user == authModel.NoneUser {
		return nil, errors.New("user not provided")
	}

	item, ok := o.(*model.Reading)
	if !ok {
		return nil, errors.New("item is not a reading")
	}

	item.User = user

	return item, nil
}
