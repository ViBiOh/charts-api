package conservatories

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/db"
	"github.com/ViBiOh/httputils/pagination"
)

const (
	healthcheckPath = `/health`
	defaultPageSize = uint(10)
	maxPageSize     = uint(50)
	defaultSort     = `name`
)

var (
	dbConfig = db.Flags(`chartsDb`)
	chartsDB *sql.DB
)

// Init charts handler
func Init() (err error) {
	chartsDB, err = db.GetDB(dbConfig)
	if err != nil {
		err = fmt.Errorf(`Error while initializing database: %v`, err)
	}

	return
}

func listCrud(w http.ResponseWriter, r *http.Request) {
	page, pageSize, sort, asc, err := pagination.ParsePaginationParams(r, defaultPageSize, maxPageSize)
	if err != nil {
		httputils.BadRequest(w, err)
		return
	}

	if sort == `` {
		sort = defaultSort
	}

	if count, list, err := findConservatories(page, pageSize, sort, asc, r.URL.Query().Get(`q`)); err != nil {
		httputils.InternalServerError(w, err)
	} else {
		httputils.ResponsePaginatedJSON(w, http.StatusOK, page, pageSize, count, list, httputils.IsPretty(r.URL.RawQuery))
	}
}

func aggregate(w http.ResponseWriter, r *http.Request) {
	if count, err := countByDepartment(); err != nil {
		httputils.InternalServerError(w, err)
	} else {
		httputils.ResponseJSON(w, 200, count, httputils.IsPretty(r.URL.RawQuery))
	}
}

func aggregateByDepartment(w http.ResponseWriter, r *http.Request, path string) {
	if count, err := countByZipOfDepartment(strings.TrimPrefix(path, `/`)); err != nil {
		httputils.InternalServerError(w, err)
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
