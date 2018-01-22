package readings

import (
	authProvider "github.com/ViBiOh/auth/provider"
)

type tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	user *authProvider.User
}
