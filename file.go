package qry

import (
	"bytes"
)

type file struct {
	*bytes.Buffer
}

func (b file) queries() QuerySet {
	var (
		raw = b.Bytes()
		str = rgxSearchQueries.FindAllSubmatch(raw, -1)
		idx = rgxSearchQueries.FindAllSubmatchIndex(raw, -1)
		out = make(QuerySet, len(str))
		key string
	)

	// Capture all from current header until next header (or EOF)
	for i := range idx {
		key = string(str[i][1])
		if len(idx) > i+1 {
			out[key] = normalize(raw[idx[i][1]:idx[i+1][0]])
		} else {
			out[key] = normalize(raw[idx[i][1]:])
		}
	}

	return out
}
