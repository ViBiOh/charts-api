package readings

import "github.com/ViBiOh/auth/pkg/model"

type reading struct {
	ID     uint   `json:"id"`
	URL    string `json:"url"`
	Public bool   `json:"public"`
	Read   bool   `json:"read"`
	Tags   []*tag `json:"tags"`
	user   *model.User
}
