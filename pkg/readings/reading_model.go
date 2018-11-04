package readings

import "github.com/ViBiOh/auth/pkg/model"

type reading struct {
	ID   uint   `json:"id"`
	URL  string `json:"url"`
	Read bool   `json:"read"`
	Tags []*tag `json:"tags"`
	user *model.User
}
