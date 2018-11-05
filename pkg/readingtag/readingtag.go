package readingtag

import (
	"database/sql"

	"github.com/ViBiOh/eponae-api/pkg/tag"
)

// App stores informations
type App struct {
	db         *sql.DB
	tagService *tag.App
}

// NewService creates new Service
func NewService(db *sql.DB, tagService *tag.App) *App {
	return &App{
		db:         db,
		tagService: tagService,
	}
}
