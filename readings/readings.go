package readings

import (
	"net/http"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils"
)

var authURL string
var authUsers map[string]*auth.User

// Init readings API
func Init(url string, users map[string]*auth.User) error {
	authURL = url
	authUsers = users

	return nil
}

// Handler for Readings request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		}

		user, err := auth.IsAuthenticated(authURL, authUsers, r)
		if err != nil {
			httputils.Unauthorized(w, err)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(user.Username))
	})
}
