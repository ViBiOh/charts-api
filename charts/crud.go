package charts

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
)

const defaultPage = int64(1)
const defaultPageSize = int64(10)
const defaultSort = `name`
const defaultOrder = true
const maxPageSize = int64(50)

func parsePaginationParams(r *http.Request) (page, pageSize int64, sortKey string, sortAsc bool, err error) {
	var parsedInt int64

	page = defaultPage
	rawPage := r.URL.Query().Get(`page`)
	if rawPage != `` {
		parsedInt, err = strconv.ParseInt(rawPage, 10, 64)
		if err != nil {
			err = fmt.Errorf(`Error while parsing page param: %v`, err)
			return
		}

		page = parsedInt
	}

	pageSize = defaultPageSize
	rawPageSize := r.URL.Query().Get(`pageSize`)
	if rawPageSize != `` {
		parsedInt, err = strconv.ParseInt(rawPageSize, 10, 64)
		if err != nil {
			err = fmt.Errorf(`Error while parsing pageSize param: %v`, err)
			return
		} else if parsedInt > maxPageSize {
			err = fmt.Errorf(`maxPageSize exceeded`)
			return
		}

		pageSize = parsedInt
	}

	sortKey = defaultSort
	rawSortKey := r.URL.Query().Get(`sort`)
	if rawSortKey != `` {
		sortKey = rawSortKey
	}

	sortAsc = defaultOrder
	rawOrder := r.URL.Query().Get(`order`)
	if rawOrder != `` {
		if rawOrder == `desc` {
			sortAsc = false
		}
	}

	return
}

func listCrud(w http.ResponseWriter, r *http.Request) {
	page, pageSize, sort, asc, err := parsePaginationParams(r)
	if err != nil {
		httputils.BadRequest(w, err)
		return
	}

	if count, list, err := findConservatories(page, pageSize, sort, asc, r.URL.Query().Get(`q`)); err != nil {
		httputils.InternalServer(w, err)
	} else {
		httputils.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list)
	}
}

func aggregate(w http.ResponseWriter, r *http.Request) {
	if count, err := countByDepartment(); err != nil {
		httputils.InternalServer(w, err)
	} else {
		httputils.ResponseJSON(w, 200, count)
	}
}

func aggregateByDepartment(w http.ResponseWriter, r *http.Request, path string) {
	if count, err := countByZipOfDepartment(strings.TrimPrefix(path, `/`)); err != nil {
		httputils.InternalServer(w, err)
	} else {
		httputils.ResponseJSON(w, 200, count)
	}
}

func conservatoriesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodGet && r.URL.Path == `/` {
			listCrud(w, r)
		} else if r.Method == http.MethodGet && r.URL.Path == `/aggregate` {
			aggregate(w, r)
		} else if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, `/aggregate/department`) {
			aggregateByDepartment(w, r, strings.TrimPrefix(r.URL.Path, `/aggregate/department`))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
