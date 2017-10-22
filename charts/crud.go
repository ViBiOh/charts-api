package charts

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
)

const defaultPage = int64(1)
const defaultPageSize = int64(20)
const defaultSort = `name`

func getRequestID(r *http.Request) (int64, error) {
	return strconv.ParseInt(strings.TrimPrefix(r.URL.Path, `/`), 10, 64)
}

func listCrud(w http.ResponseWriter, r *http.Request) {
	page := defaultPage
	rawPage := r.URL.Query().Get(`page`)
	if rawPage != `` {
		parsedPage, err := strconv.ParseInt(rawPage, 10, 64)
		if err != nil {
			httputils.BadRequest(w, fmt.Errorf(`Error while parsing page param: %v`, err))
			return
		}

		page = parsedPage
	}

	pageSize := defaultPageSize
	rawPageSize := r.URL.Query().Get(`pageSize`)
	if rawPageSize != `` {
		parsedPageSize, err := strconv.ParseInt(rawPageSize, 10, 64)
		if err != nil {
			httputils.BadRequest(w, fmt.Errorf(`Error while parsing pageSize param: %v`, err))
			return
		}

		pageSize = parsedPageSize
	}

	sortKey := defaultSort
	rawSortKey := r.URL.Query().Get(`o`)
	if rawSortKey != `` {
		sortKey = rawSortKey
	}

	if list, err := listConservatories(page, pageSize, sortKey); err != nil {
		httputils.InternalServer(w, err)
	} else {
		httputils.ResponseArrayJSON(w, http.StatusOK, list)
	}
}

func conservatoriesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodGet && r.URL.Path == `/` {
			listCrud(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
