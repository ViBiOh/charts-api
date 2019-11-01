package model

import (
	authModel "github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
)

var _ crud.Item = &Reading{}

// Reading describe an url saved by an user
type Reading struct {
	ID   string          `json:"id"`
	URL  string          `json:"url"`
	Read bool            `json:"read"`
	Tags []*Tag          `json:"tags"`
	User *authModel.User `json:"-"`
}

// SetID setter for ID
func (o *Reading) SetID(id string) {
	o.ID = id
}
