package model

import (
	authModel "github.com/ViBiOh/auth/pkg/model"
)

// Tag describe a meta label defined by an user
type Tag struct {
	UUID string          `json:"id"`
	Name string          `json:"name"`
	User *authModel.User `json:"-"`
}

// ID returns UUID
func (a Tag) ID() string {
	return a.UUID
}
