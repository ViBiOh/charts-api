package readings

import (
	"github.com/ViBiOh/auth/auth"
)

type tag struct {
	id   int64
	Name string     `json:"name"`
	User *auth.User `json:"user"`
}

type reading struct {
	id     int64
	URL    string `json:"url"`
	Public bool   `json:"public"`
	Read   bool   `json:"read"`
	Tags   []*tag `json:"tags"`
}
