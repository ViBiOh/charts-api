package readings

import (
	"database/sql"
	"encoding/json"
	native_errors "errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/auth/pkg/model"
	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/pagination"
	"github.com/ViBiOh/httputils/pkg/request"
)

const (
	tagsPath        = `/tags`
	defaultPageSize = 50
	maxPageSize     = uint(^uint(0) >> 1)
)

var (
	errNameRequired = native_errors.New(`name is required`)
	errIDRequired   = native_errors.New(`iD is required in path`)
)

func getRequestID(path string) (uint, error) {
	parsed, err := strconv.ParseUint(strings.TrimPrefix(path, `/`), 10, 32)
	return uint(parsed), err
}

func (a App) readTagFromBody(r *http.Request) (*tag, error) {
	var requestTag tag

	if bodyBytes, err := request.ReadBodyRequest(r); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bodyBytes, &requestTag); err != nil {
		return nil, errors.WithStack(err)
	}

	return &requestTag, nil
}

func (a App) listTags(w http.ResponseWriter, r *http.Request, user *model.User) {
	page, pageSize, sort, asc, err := pagination.ParseParams(r, 1, defaultPageSize, maxPageSize)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if sort == `` {
		sort = `name`
	}

	query := r.URL.Query().Get(`q`)

	list, err := a.searchTags(page, pageSize, sort, asc, user, query)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	count, err := a.countTags(user, query)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := httpjson.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) readTag(w http.ResponseWriter, r *http.Request, user *model.User, id uint) {
	foundTag, err := a.getTag(id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			httperror.NotFound(w)
		} else {
			httperror.InternalServerError(w, err)
		}
		return
	}

	if err := httpjson.ResponseJSON(w, http.StatusOK, foundTag, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) createTag(w http.ResponseWriter, r *http.Request, user *model.User) {
	bodyTag, err := a.readTagFromBody(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if bodyTag.Name == `` {
		httperror.BadRequest(w, errNameRequired)
		return
	}

	bodyTag.user = user

	if err := a.saveTag(bodyTag, nil); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := httpjson.ResponseJSON(w, http.StatusCreated, bodyTag, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) updateTag(w http.ResponseWriter, r *http.Request, user *model.User, id uint) {
	bodyTag, err := a.readTagFromBody(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if id == 0 {
		httperror.BadRequest(w, errIDRequired)
		return
	}

	if bodyTag.Name == `` {
		httperror.BadRequest(w, errNameRequired)
		return
	}

	bodyTag.ID = id
	bodyTag.user = user

	if err := a.saveTag(bodyTag, nil); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := httpjson.ResponseJSON(w, http.StatusCreated, bodyTag, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) removeTag(w http.ResponseWriter, r *http.Request, user *model.User, id uint) {
	foundTag, err := a.getTag(id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			httperror.NotFound(w)
		} else {
			httperror.InternalServerError(w, err)
		}
		return
	}

	if err := a.deleteTag(foundTag, nil); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		id, err := getRequestID(path)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}

		if r.Method == http.MethodGet {
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
