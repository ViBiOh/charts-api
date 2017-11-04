package readings

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils"
)

func createTag(w http.ResponseWriter, r *http.Request, user *auth.User) {
	var requestTag tag

	if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while reading body: %v`, err))
	} else if err := json.Unmarshal(bodyBytes, &requestTag); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while unmarshalling body: %v`, err))
	} else if requestTag.Name == `` {
		httputils.BadRequest(w, errors.New(`Name is required`))
	} else {
		requestTag.user = user

		if err := saveTag(&requestTag, nil); err != nil {
			httputils.InternalServerError(w, err)
		} else {
			httputils.ResponseJSON(w, http.StatusCreated, requestTag, httputils.IsPretty(r.URL.RawQuery))
		}
	}
}
