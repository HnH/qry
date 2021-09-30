package qry

import (
	"bytes"
	"strings"
)

// QuerySet represents set of queries
type QuerySet map[string]Query

// Query represents single query from a .sql file
type Query string

// Replace part of a query
func (q Query) Replace(o, r string) Query {
	if len(o) == 0 || len(r) == 0 {
		return q
	}

	return Query(strings.Replace(string(q), o, r, 1))
}

func removeMultilineComments(q []byte) []byte {
	for {
		index := bytes.IndexAny(q, "/*")
		if index == -1 {
			return q
		}

		closeIndex := bytes.Index(q, []byte("*/"))
		if closeIndex == -1 {
			return q
		}

		q = append(q[:index], q[closeIndex+2:]...)
	}
}

func normalize(q []byte) Query {
	q = bytes.TrimSpace(q)
	q = removeMultilineComments(q)
	q = rgxLineComment.ReplaceAll(q, []byte(" "))
	q = bytes.Replace(q, []byte("\n"), []byte(" "), -1)
	q = rgxMultiSpace.ReplaceAll(q, []byte(" "))
	q = bytes.Replace(q, []byte("\\"), []byte("\\\\"), -1)
	q = bytes.Replace(q, []byte("\""), []byte("\\\""), -1)

	return Query(q)
}

// In returns string with N sql query placeholders
func In(l int) string {
	if l <= 0 {
		return ""
	}

	return strings.Repeat(ph, l)[1:]
}
