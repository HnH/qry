package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/HnH/qry"
)

type config struct {
	dir string
	pkg string
	out string
	fmt bool
}

func main() {
	var cfg, err = loadCfg()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err = gen(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func loadCfg() (cfg config, err error) {
	flag.StringVar(&cfg.dir, "dir", "", "root directory with .sql files")
	flag.StringVar(&cfg.pkg, "pkg", "", "go package directory")
	flag.StringVar(&cfg.out, "out", "qry.go", "output filename")
	flag.BoolVar(&cfg.fmt, "fmt", true, "pass output trough gofmt")
	flag.Parse()

	if len(cfg.dir) == 0 {
		if cfg.dir, err = os.Getwd(); err != nil {
			err = qry.ErrDirSql
		}
	}

	return
}

func gen(cfg config) (err error) {
	var b *bytes.Buffer
	if b, err = loadSql(cfg); err != nil {
		return
	}

	// passing trough gofmt
	if cfg.fmt {
		if err = format(b); err != nil {
			return
		}
	}

	// Output to stdout
	if len(cfg.pkg) == 0 {
		fmt.Fprintf(os.Stdout, b.String())
		return
	}

	// Output to file
	if err = ioutil.WriteFile(fmt.Sprintf("%s/%s", cfg.pkg, cfg.out), b.Bytes(), os.ModePerm); err != nil {
		err = fmt.Errorf("cannot write to output file [%s/%s]", cfg.pkg, cfg.out)
		return
	}

	fmt.Fprintf(os.Stdout, "saved to output file [%s/%s]\n", cfg.pkg, cfg.out)
	return
}

func loadSql(cfg config) (b *bytes.Buffer, err error) {
	var loaded []qry.File
	if loaded, err = qry.DirOrdered(cfg.dir); err != nil {
		err = qry.ErrDirSql
		return
	}

	b = bytes.NewBuffer(nil)

	// Package directory provided
	if len(cfg.pkg) > 0 {
		var stat os.FileInfo
		if stat, err = os.Stat(cfg.pkg); err != nil || !stat.IsDir() {
			err = qry.ErrDirPkg
			return
		}

		b.WriteString(fmt.Sprintf("package %s\n\n", filepath.Base(cfg.pkg)))
	}

	b.WriteString("const (\n")

	for _, f := range loaded {
		b.WriteString(fmt.Sprintf("// %s\n", f.Name))

		for _, i := range f.Items {
			b.WriteString(fmt.Sprintf("\t%s = \"%s\"\n", i.Name, i.Query))
		}

		b.WriteString("\n\n")
	}

	b.WriteString(")")

	return
}

func format(b *bytes.Buffer) (err error) {
	var (
		cmd   = exec.Command("gofmt") // #nosec
		stdin io.WriteCloser
	)

	if stdin, err = cmd.StdinPipe(); err != nil {
		err = errors.New("cannot acquire gofmt stdin")
		return
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, b.String())
	}()

	var out []byte
	if out, err = cmd.CombinedOutput(); err != nil {
		err = errors.New("cannot acquire gofmt stdout")
		return
	}

	b.Reset()
	b.Write(out)

	return
}
