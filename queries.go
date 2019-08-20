package qry

import (
	"strings"
)

// QuerySet represents set of queries
type QuerySet map[string]Query

// Query represents single query from a .sql file
type Query string

// Replaces part of a query
func (q Query) Replace(o, r string) Query {
	if len(o) == 0 || len(r) == 0 {
		return q
	}

	return Query(strings.Replace(string(q), o, r, 1))
}

// In returns string with N sql query placeholders
func In(l int) string {
	if l <= 0 {
		return ""
	}

	return strings.Repeat(ph, l)[1:]
}
