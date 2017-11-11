package readings

import (
	"github.com/ViBiOh/auth/auth"
)

type tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	user *auth.User
}
