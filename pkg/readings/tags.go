package readings

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/pagination"
	"github.com/ViBiOh/httputils/pkg/request"
)

const tagsPath = `/tags`
const defaultPageSize = 50
const maxPageSize = uint(^uint(0) >> 1)

var errNameRequired = errors.New(`name is required`)
var errIDRequired = errors.New(`iD is required in path`)

func getRequestID(path string) (uint, error) {
	parsed, err := strconv.ParseUint(strings.TrimPrefix(path, `/`), 10, 32)
	return uint(parsed), err
}

func (a App) readTagFromBody(r *http.Request) (*tag, error) {
	var requestTag tag

	if bodyBytes, err := request.ReadBodyRequest(r); err != nil {
		return nil, fmt.Errorf(`error while reading body: %v`, err)
	} else if err := json.Unmarshal(bodyBytes, &requestTag); err != nil {
		return nil, fmt.Errorf(`error while unmarshalling body: %v`, err)
	}

	return &requestTag, nil
}

func (a App) listTags(w http.ResponseWriter, r *http.Request, user *model.User) {
	page, pageSize, sort, asc, err := pagination.ParseParams(r, 1, defaultPageSize, maxPageSize)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf(`error while parsing pagination: %v`, err))
		return
	}

	if sort == `` {
		sort = `name`
	}

	query := r.URL.Query().Get(`q`)

	if list, err := a.searchTags(page, pageSize, sort, asc, user, query); err != nil {
		httperror.InternalServerError(w, fmt.Errorf(`error while searching tags: %v`, err))
	} else if count, err := a.countTags(user, query); err != nil {
		httperror.InternalServerError(w, fmt.Errorf(`error while counting tags: %v`, err))
	} else if err := httpjson.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func (a App) readTag(w http.ResponseWriter, r *http.Request, user *model.User, id uint) {
	if foundTag, err := a.getTag(id, user); err != nil {
		if err == sql.ErrNoRows {
			httperror.NotFound(w)
		} else {
			httperror.InternalServerError(w, fmt.Errorf(`error while getting tag: %v`, err))
		}
	} else if err := httpjson.ResponseJSON(w, http.StatusOK, foundTag, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func (a App) createTag(w http.ResponseWriter, r *http.Request, user *model.User) {
	if bodyTag, err := a.readTagFromBody(r); err != nil {
		httperror.BadRequest(w, fmt.Errorf(`error while parsing body: %v`, err))
	} else if bodyTag.Name == `` {
		httperror.BadRequest(w, errNameRequired)
	} else {
		bodyTag.user = user

		if err := a.saveTag(bodyTag, nil); err != nil {
			httperror.InternalServerError(w, fmt.Errorf(`error while saving tag: %v`, err))
		} else if err := httpjson.ResponseJSON(w, http.StatusCreated, bodyTag, httpjson.IsPretty(r)); err != nil {
			httperror.InternalServerError(w, err)
		}
	}
}
func (a App) updateTag(w http.ResponseWriter, r *http.Request, user *model.User, id uint) {
	if bodyTag, err := a.readTagFromBody(r); err != nil {
		httperror.BadRequest(w, fmt.Errorf(`error while parsing body: %v`, err))
	} else if id == 0 {
		httperror.BadRequest(w, errIDRequired)
	} else if bodyTag.Name == `` {
		httperror.BadRequest(w, errNameRequired)
	} else {
		bodyTag.ID = id
		bodyTag.user = user

		if err := a.saveTag(bodyTag, nil); err != nil {
			httperror.InternalServerError(w, fmt.Errorf(`error while saving tag: %v`, err))
		} else if err := httpjson.ResponseJSON(w, http.StatusCreated, bodyTag, httpjson.IsPretty(r)); err != nil {
			httperror.InternalServerError(w, err)
		}
	}
}

func (a App) removeTag(w http.ResponseWriter, r *http.Request, user *model.User, id uint) {
	if foundTag, err := a.getTag(id, user); err != nil {
		if err == sql.ErrNoRows {
			httperror.NotFound(w)
		} else {
			httperror.InternalServerError(w, fmt.Errorf(`error while getting tag: %v`, err))
		}
	} else if err := a.deleteTag(foundTag, nil); err != nil {
		httperror.InternalServerError(w, fmt.Errorf(`error while deleting tag: %v`, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a App) tagsHandler(w http.ResponseWriter, r *http.Request, user *model.User, path string) {
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
			httperror.BadRequest(w, fmt.Errorf(`error while parsing request path: %v`, err))
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
