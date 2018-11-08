package model

import (
	authModel "github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/crud"
)

var _ crud.Item = &Tag{}

// Tag describe a meta label defined by an user
type Tag struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	User *authModel.User `json:"-"`
}

// SetID setter for ID
func (o *Tag) SetID(id string) {
	o.ID = id
}
