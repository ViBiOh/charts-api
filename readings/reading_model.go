package readings

import (
	authProvider "github.com/ViBiOh/auth/provider"
)

type reading struct {
	ID     uint   `json:"id"`
	URL    string `json:"url"`
	Public bool   `json:"public"`
	Read   bool   `json:"read"`
	Tags   []*tag `json:"tags"`
	user   *authProvider.User
}
