package model

import (
	authModel "github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/crud"
)

var _ crud.Item = Reading{}

// Reading describe an url saved by an user
type Reading struct {
	UUID string          `json:"id"`
	URL  string          `json:"url"`
	Read bool            `json:"read"`
	Tags []*Tag          `json:"tags"`
	User *authModel.User `json:"-"`
}

// ID returns UUID
func (a Reading) ID() string {
	return a.UUID
}
