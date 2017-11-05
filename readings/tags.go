package readings

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/auth/auth"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/pagination"
)

const tagsPath = `/tags`
const defaultPageSize = 50
const maxPageSize = int64(^uint(0) >> 1)

func getRequestID(path string) (int64, error) {
	return strconv.ParseInt(strings.TrimPrefix(path, `/`), 10, 64)
}

func readTagFromBody(r *http.Request) (*tag, error) {
	var requestTag tag

	if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
		return nil, fmt.Errorf(`Error while reading body: %v`, err)
	} else if err := json.Unmarshal(bodyBytes, &requestTag); err != nil {
		return nil, fmt.Errorf(`Error while unmarshalling body: %v`, err)
	}

	return &requestTag, nil
}

func listTags(w http.ResponseWriter, r *http.Request, user *auth.User) {
	page, pageSize, sort, asc, err := pagination.ParsePaginationParams(r, defaultPageSize, maxPageSize)
	if err != nil {
		httputils.BadRequest(w, err)
		return
	}

	if sort == `` {
		sort = `name`
	}

	query := r.URL.Query().Get(`q`)

	if list, err := searchTags(page, pageSize, sort, asc, user, query); err != nil {
		httputils.InternalServerError(w, err)
	} else if count, err := countTags(user, query); err != nil {
		httputils.InternalServerError(w, err)
	} else {
		httputils.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list, httputils.IsPretty(r.URL.RawQuery))
	}
}

func readTag(w http.ResponseWriter, r *http.Request, user *auth.User, id int64) {
	if foundTag, err := getTag(id, user); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			httputils.InternalServerError(w, fmt.Errorf(`Error while getting tag: %v`, err))
		}
	} else {
		httputils.ResponseJSON(w, http.StatusOK, foundTag, httputils.IsPretty(r.URL.RawQuery))
	}
}

func createTag(w http.ResponseWriter, r *http.Request, user *auth.User) {
	if bodyTag, err := readTagFromBody(r); err != nil {
		httputils.BadRequest(w, err)
	} else if bodyTag.Name == `` {
		httputils.BadRequest(w, errors.New(`Name is required`))
	} else {
		bodyTag.user = user

		if err := saveTag(bodyTag, nil); err != nil {
			httputils.InternalServerError(w, err)
		} else {
			httputils.ResponseJSON(w, http.StatusCreated, bodyTag, httputils.IsPretty(r.URL.RawQuery))
		}
	}
}
func updateTag(w http.ResponseWriter, r *http.Request, user *auth.User, id int64) {
	if bodyTag, err := readTagFromBody(r); err != nil {
		httputils.BadRequest(w, err)
	} else if id == 0 {
		httputils.BadRequest(w, errors.New(`ID is required in path`))
	} else if bodyTag.Name == `` {
		httputils.BadRequest(w, errors.New(`Name is required`))
	} else {
		bodyTag.ID = id
		bodyTag.user = user

		if err := saveTag(bodyTag, nil); err != nil {
			httputils.InternalServerError(w, err)
		} else {
			httputils.ResponseJSON(w, http.StatusCreated, bodyTag, httputils.IsPretty(r.URL.RawQuery))
		}
	}
}

func removeTag(w http.ResponseWriter, r *http.Request, user *auth.User, id int64) {
	if foundTag, err := getTag(id, user); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			httputils.InternalServerError(w, fmt.Errorf(`Error while getting tag for deletion: %v`, err))
		}
	} else if err := deleteTag(foundTag, nil); err != nil {
		httputils.InternalServerError(w, fmt.Errorf(`Error while deleting tag: %v`, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func tagsHandler(w http.ResponseWriter, r *http.Request, user *auth.User, path string) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	} else if path == `/` || path == `` {
		if r.Method == http.MethodPost {
			createTag(w, r, user)
		} else if r.Method == http.MethodGet {
			listTags(w, r, user)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		if id, err := getRequestID(path); err != nil {
			httputils.BadRequest(w, err)
		} else if r.Method == http.MethodGet {
			readTag(w, r, user, id)
		} else if r.Method == http.MethodPut {
			updateTag(w, r, user, id)
		} else if r.Method == http.MethodDelete {
			removeTag(w, r, user, id)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
