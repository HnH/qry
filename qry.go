/*
Package qry is a general purpose library for storing your raw database queries in .sql files with all benefits of modern IDEs,
instead of strings and constants in the code, and using them in an easy way inside your application with all the profit
of compile time constants.

qry recursively loads all .sql files from a specified folder, parses them according to predefined rules and returns a
reusable object, which is actually just a `map[string]string` with some sugar. Multiple queries inside a single file
are separated with standard SQL comment syntax: `-- qry: QueryName`. A `QueryName` must match `[A-Za-z_]+`.

gen tool is used for automatic generation of constants for all user specified `query_names`.
*/
package qry

import (
	"errors"
	"os"
	"regexp"
)

const (
	ext = ".sql"
	ph  = ",?"
)

var (
	rgxSearchQueries = regexp.MustCompile(`(?m)^--\s*qry:\s*([A-Za-z_]+)\s*$`)
	rgxMultiSpace    = regexp.MustCompile(`\s{2,}`)
	rgxLineComment   = regexp.MustCompile(`--[^\n]*`)

	// ErrDirSql is returned in case when directory with .sql files is unavailable
	ErrDirSql = errors.New("cannot find directory with .sql files")
	// ErrDirPkg is returned in case when directory with go package is unavailable
	ErrDirPkg = errors.New("cannot find go package directory")
)

// Dir recursively loads all .sql files from a specified folder and returns them as a hashmap
func Dir(dir string) (queries map[string]QuerySet, err error) {
	var files []File
	if files, err = DirOrdered(dir); err != nil {
		return
	}

	queries = make(map[string]QuerySet, len(files))
	for _, f := range files {
		queries[f.Name] = make(QuerySet, len(f.Items))

		for _, i := range f.Items {
			queries[f.Name][i.Name] = i.Query
		}
	}

	return queries, nil
}

// DirOrdered recursively loads all .sql files from a specified folder and returns them as a slice
func DirOrdered(dir string) ([]File, error) {
	if s, err := os.Stat(dir); err != nil || !s.IsDir() {
		return nil, ErrDirSql
	}

	return readFiles(dir)
}
