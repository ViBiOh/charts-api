package readingtag

import (
	"database/sql"

	"github.com/ViBiOh/eponae-api/pkg/tag"
)

// App of package
type App struct {
	db         *sql.DB
	tagService *tag.App
}

// New creates new App
func New(db *sql.DB, tagService *tag.App) *App {
	return &App{
		db:         db,
		tagService: tagService,
	}
}
