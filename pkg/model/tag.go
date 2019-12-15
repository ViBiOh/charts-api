package model

import (
	authModel "github.com/ViBiOh/auth/v2/pkg/model"
)

// Tag describe a meta label defined by an user
type Tag struct {
	ID   uint64         `json:"id"`
	Name string         `json:"name"`
	User authModel.User `json:"-"`
}
