package charts

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/db"
)

const healthcheckPath = `/health`
const defaultPage = int64(1)
const defaultPageSize = int64(10)
const defaultSort = `name`
const defaultOrder = true
const maxPageSize = int64(50)

var chartsDB *sql.DB
var dbConfig = db.Flags(`db`)

// Init charts handler
func Init() (err error) {
	chartsDB, err = db.GetDB(dbConfig)
	if err != nil {
		log.Printf(`[charts] Error while initializing database: %v`, err)
	} else if chartsDB != nil {
		log.Print(`[charts] Database ready`)
	}

	return nil
}

func parsePaginationParams(r *http.Request) (page, pageSize int64, sortKey string, sortAsc bool, err error) {
	var parsedInt int64
	var params url.Values

	params, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return
	}

	page = defaultPage
	rawPage := params.Get(`page`)
	if rawPage != `` {
		parsedInt, err = strconv.ParseInt(rawPage, 10, 64)
		if err != nil {
			err = fmt.Errorf(`Error while parsing page param: %v`, err)
			return
		}

		page = parsedInt
	}

	pageSize = defaultPageSize
	rawPageSize := params.Get(`pageSize`)
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
	rawSortKey := params.Get(`sort`)
	if rawSortKey != `` {
		sortKey = rawSortKey
	}

	sortAsc = defaultOrder
	if _, ok := params[`desc`]; ok {
		sortAsc = false
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
		httputils.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list, httputils.IsPretty(r.URL.RawQuery))
	}
}

func aggregate(w http.ResponseWriter, r *http.Request) {
	if count, err := countByDepartment(); err != nil {
		httputils.InternalServer(w, err)
	} else {
		httputils.ResponseJSON(w, 200, count, httputils.IsPretty(r.URL.RawQuery))
	}
}

func aggregateByDepartment(w http.ResponseWriter, r *http.Request, path string) {
	if count, err := countByZipOfDepartment(strings.TrimPrefix(path, `/`)); err != nil {
		httputils.InternalServer(w, err)
	} else {
		httputils.ResponseJSON(w, 200, count, httputils.IsPretty(r.URL.RawQuery))
	}
}

// Handler for Charts request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == healthcheckPath {
			if db.Ping(chartsDB) {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			return
		}

		if r.Method == http.MethodGet && (r.URL.Path == `/` || r.URL.Path == ``) {
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
