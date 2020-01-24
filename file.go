package qry

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// File represents parsed file with .sql queries
type File struct {
	Name  string
	Items []FileItem
}

// FileItem represents single query within a file
type FileItem struct {
	Name  string
	Query Query
}

func readFiles(dir string) (out []File, err error) {
	var (
		buf = bytes.NewBuffer(nil)
		one FileItem
	)

	err = filepath.Walk(dir, func(path string, finfo os.FileInfo, err error) error {
		if finfo.IsDir() || filepath.Ext(finfo.Name()) != ext {
			return nil
		}

		var f *os.File
		if f, err = os.Open(filepath.Clean(path)); err != nil {
			return err
		}

		// write file contents to buffer
		buf.Reset()
		buf.WriteString("\n")
		io.Copy(buf, f)
		f.Close()

		if buf.Len() == 0 {
			return nil
		}

		// look for queries
		var (
			raw = buf.Bytes()
			str = rgxSearchQueries.FindAllSubmatch(raw, -1)
			idx = rgxSearchQueries.FindAllSubmatchIndex(raw, -1)
			fl  = File{
				Name:  finfo.Name(),
				Items: make([]FileItem, 0, len(str)),
			}
		)

		// Capture all from current header until next header (or EOF)
		for i := range idx {
			one.Name = string(str[i][1])
			if len(idx) > i+1 {
				one.Query = normalize(raw[idx[i][1]:idx[i+1][0]])
			} else {
				one.Query = normalize(raw[idx[i][1]:])
			}

			fl.Items = append(fl.Items, one)
			one = FileItem{}
		}

		out = append(out, fl)

		return nil
	})

	return
}
