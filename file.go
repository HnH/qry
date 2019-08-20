package qry

import (
	"bytes"
	"strings"
)

type file struct {
	*bytes.Buffer
}

func (b file) queries() QuerySet {
	var (
		raw = b.String()
		str = rgxSearchQueries.FindAllStringSubmatch(raw, -1)
		idx = rgxSearchQueries.FindAllStringSubmatchIndex(raw, -1)
		out = make(QuerySet, len(str))
		key string
	)

	// Capture all from current header until next header (or EOF)
	for i := range idx {
		key = str[i][1]
		if len(idx) > i+1 {
			out[key] = Query(strings.TrimSpace(strings.Replace(raw[idx[i][1]:idx[i+1][0]], "\n", "", -1)))
		} else {
			out[key] = Query(strings.TrimSpace(strings.Replace(raw[idx[i][1]:], "\n", "", -1)))
		}
	}

	return out
}
