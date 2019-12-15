package model

import (
	authModel "github.com/ViBiOh/auth/v2/pkg/model"
)

// Reading describe an url saved by an user
type Reading struct {
	ID   uint64         `json:"id"`
	URL  string         `json:"url"`
	Read bool           `json:"read"`
	Tags []*Tag         `json:"tags"`
	User authModel.User `json:"-"`
}
