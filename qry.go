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
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

const (
	ext = ".sql"
	ph  = ",?"
)

var (
	rgxSearchQueries = regexp.MustCompile(`(?m)^--\s*qry:\s*([A-Za-z_]+)\s*$`)

	ErrDirSql = errors.New("cannot find directory with .sql files")
	ErrDirPkg = errors.New("cannot find go package directory")
)

// Dir recursively loads all .sql files from a specified folder
func Dir(dir string) (queries map[string]QuerySet, err error) {
	var s os.FileInfo
	if s, err = os.Stat(dir); err != nil || !s.IsDir() {
		err = ErrDirSql
		return
	}

	var files map[string]file
	if files, err = readFiles(dir); err != nil {
		return
	}

	queries = make(map[string]QuerySet)
	for filename, buffer := range files {
		if buffer.Len() == 0 {
			continue
		}

		queries[filename] = buffer.queries()
	}

	return queries, nil
}

func readFiles(dir string) (map[string]file, error) {
	var out = make(map[string]file)
	var err = filepath.Walk(dir, func(path string, finfo os.FileInfo, err error) error {
		if finfo.IsDir() || filepath.Ext(finfo.Name()) != ext {
			return nil
		}

		var f *os.File
		if f, err = os.Open(filepath.Clean(path)); err != nil {
			return err
		}

		if _, ok := out[finfo.Name()]; !ok {
			out[finfo.Name()] = file{
				bytes.NewBuffer(nil),
			}
		}

		var b = out[finfo.Name()]
		b.WriteString("\n")
		io.Copy(b, f)
		f.Close()

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
