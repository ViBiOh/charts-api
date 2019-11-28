package model

import (
	authModel "github.com/ViBiOh/auth/v2/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
)

var _ crud.Item = &Reading{}

// Reading describe an url saved by an user
type Reading struct {
	ID   uint64         `json:"id"`
	URL  string         `json:"url"`
	Read bool           `json:"read"`
	Tags []*Tag         `json:"tags"`
	User authModel.User `json:"-"`
}

// SetID setter for ID
func (o *Reading) SetID(id uint64) {
	o.ID = id
}
