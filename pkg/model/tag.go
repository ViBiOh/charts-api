package model

import (
	authModel "github.com/ViBiOh/auth/v2/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
)

var _ crud.Item = &Tag{}

// Tag describe a meta label defined by an user
type Tag struct {
	ID   uint64         `json:"id"`
	Name string         `json:"name"`
	User authModel.User `json:"-"`
}

// SetID setter for ID
func (o *Tag) SetID(id uint64) {
	o.ID = id
}
