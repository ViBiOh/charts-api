package readings

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	authProvider "github.com/ViBiOh/auth/provider"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/pagination"
)

const tagsPath = `/tags`
const defaultPageSize = 50
const maxPageSize = uint(^uint(0) >> 1)

var errNameRequired = errors.New(`Name is required`)
var errIDRequired = errors.New(`ID is required in path`)

func getRequestID(path string) (uint, error) {
	parsed, err := strconv.ParseUint(strings.TrimPrefix(path, `/`), 10, 32)
	return uint(parsed), err
}

func (a *App) readTagFromBody(r *http.Request) (*tag, error) {
	var requestTag tag

	if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
		return nil, fmt.Errorf(`Error while reading body: %v`, err)
	} else if err := json.Unmarshal(bodyBytes, &requestTag); err != nil {
		return nil, fmt.Errorf(`Error while unmarshalling body: %v`, err)
	}

	return &requestTag, nil
}

func (a *App) listTags(w http.ResponseWriter, r *http.Request, user *authProvider.User) {
	page, pageSize, sort, asc, err := pagination.ParsePaginationParams(r, defaultPageSize, maxPageSize)
	if err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing pagination: %v`, err))
		return
	}

	if sort == `` {
		sort = `name`
	}

	query := r.URL.Query().Get(`q`)

	if list, err := a.searchTags(page, pageSize, sort, asc, user, query); err != nil {
		httputils.InternalServerError(w, fmt.Errorf(`Error while searching tags: %v`, err))
	} else if count, err := a.countTags(user, query); err != nil {
		httputils.InternalServerError(w, fmt.Errorf(`Error while counting tags: %v`, err))
	} else {
		httputils.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list, httputils.IsPretty(r.URL.RawQuery))
	}
}

func (a *App) readTag(w http.ResponseWriter, r *http.Request, user *authProvider.User, id uint) {
	if foundTag, err := a.getTag(id, user); err != nil {
		if err == sql.ErrNoRows {
			httputils.NotFound(w)
		} else {
			httputils.InternalServerError(w, fmt.Errorf(`Error while getting tag: %v`, err))
		}
	} else {
		httputils.ResponseJSON(w, http.StatusOK, foundTag, httputils.IsPretty(r.URL.RawQuery))
	}
}

func (a *App) createTag(w http.ResponseWriter, r *http.Request, user *authProvider.User) {
	if bodyTag, err := a.readTagFromBody(r); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing body: %v`, err))
	} else if bodyTag.Name == `` {
		httputils.BadRequest(w, errNameRequired)
	} else {
		bodyTag.user = user

		if err := a.saveTag(bodyTag, nil); err != nil {
			httputils.InternalServerError(w, fmt.Errorf(`Error while saving tag: %v`, err))
		} else {
			httputils.ResponseJSON(w, http.StatusCreated, bodyTag, httputils.IsPretty(r.URL.RawQuery))
		}
	}
}
func (a *App) updateTag(w http.ResponseWriter, r *http.Request, user *authProvider.User, id uint) {
	if bodyTag, err := a.readTagFromBody(r); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing body: %v`, err))
	} else if id == 0 {
		httputils.BadRequest(w, errIDRequired)
	} else if bodyTag.Name == `` {
		httputils.BadRequest(w, errNameRequired)
	} else {
		bodyTag.ID = id
		bodyTag.user = user

		if err := a.saveTag(bodyTag, nil); err != nil {
			httputils.InternalServerError(w, fmt.Errorf(`Error while saving tag: %v`, err))
		} else {
			httputils.ResponseJSON(w, http.StatusCreated, bodyTag, httputils.IsPretty(r.URL.RawQuery))
		}
	}
}

func (a *App) removeTag(w http.ResponseWriter, r *http.Request, user *authProvider.User, id uint) {
	if foundTag, err := a.getTag(id, user); err != nil {
		if err == sql.ErrNoRows {
			httputils.NotFound(w)
		} else {
			httputils.InternalServerError(w, fmt.Errorf(`Error while getting tag: %v`, err))
		}
	} else if err := a.deleteTag(foundTag, nil); err != nil {
		httputils.InternalServerError(w, fmt.Errorf(`Error while deleting tag: %v`, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a *App) tagsHandler(w http.ResponseWriter, r *http.Request, user *authProvider.User, path string) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	} else if path == `/` || path == `` {
		if r.Method == http.MethodPost {
			a.createTag(w, r, user)
		} else if r.Method == http.MethodGet {
			a.listTags(w, r, user)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		if id, err := getRequestID(path); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`Error while parsing request path: %v`, err))
		} else if r.Method == http.MethodGet {
			a.readTag(w, r, user, id)
		} else if r.Method == http.MethodPut {
			a.updateTag(w, r, user, id)
		} else if r.Method == http.MethodDelete {
			a.removeTag(w, r, user, id)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
