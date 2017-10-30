package readings

import (
	"net/http"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils"
)

const healthcheckPath = `/health`

var authURL string
var authUsers map[string]*auth.User

var authConfig = auth.Flags(`readingsAuth`)

// Init readings API
func Init() error {
	authURL = *authConfig[`url`]
	authUsers = auth.LoadUsersProfiles(*authConfig[`users`])

	return nil
}

// Handler for Readings request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == healthcheckPath {
			w.WriteHeader(http.StatusOK)
			return
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
