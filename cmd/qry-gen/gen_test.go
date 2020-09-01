package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestGen(t *testing.T) {
	var (
		cfg = config{
			dir:     "../../testdata",
			out:     "gen.go",
			fmt:     true,
			comment: true,
		}

		err error
	)

	cfg.pkg = cfg.dir + "/pkg/"
	os.RemoveAll(cfg.pkg)
	defer os.RemoveAll(cfg.pkg)

	if err = os.Mkdir(cfg.pkg, 0750); err != nil {
		t.Errorf("could not create pkg dir: %s", err.Error())
	}

	if err = gen(cfg); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	b, err := ioutil.ReadFile(cfg.pkg + cfg.out)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	//revive:disable:line-length-limit
	for i, c := range []bool{
		bytes.Contains(b, []byte("package pkg")),
		bytes.Contains(b, []byte("\t// one.sql\n\n")),
		bytes.Contains(b, []byte("\t// InsertUser query")),
		bytes.Contains(b, []byte("\tInsertUser = \"INSERT INTO `users` (`name`) VALUES (?);\"")),
		bytes.Contains(b, []byte("\t// GetUserById query")),
		bytes.Contains(b, []byte("\tGetUserById = \"SELECT * FROM `users` WHERE `user_id` = ?;\"")),
		bytes.Contains(b, []byte("\t// two.sql\n\n")),
		bytes.Contains(b, []byte("\t// DeleteUsersByIds query")),
		bytes.Contains(b, []byte("\tDeleteUsersByIds = \"DELETE FROM `users` WHERE `user_id` IN ({ids});\"")),
		bytes.Contains(b, []byte("\t// UglyMultiLineQuery query")),
		bytes.Contains(b, []byte("\tUglyMultiLineQuery = \"SELECT * FROM `users` WHERE YEAR(`birth_date`) > 2000;\"")),
		bytes.Contains(b, []byte("\t// EscapedJSONQuery query")),
		bytes.Contains(b, []byte("\tEscapedJSONQuery = \"INSERT INTO \\\"data\\\" (id, \\\"data\\\") VALUES (1, '{\\\"test\\\": 1}'), (2, '{\\\"test\\\": 2}');\"")),
		bytes.Contains(b, []byte("\t// EscapedByteaQuery query")),
		bytes.Contains(b, []byte("\tEscapedByteaQuery = \"INSERT INTO bin (id, \\\"data\\\") VALUES (1, E'\\\\\\\\x3aaab6e7fb7245cc9785653a0d9ffc4a5ce0f974');\"")),
	} {
		if !c {
			t.Errorf("check at idx %d has not passed", i)
			t.Logf("%s", b)
		}
	}
	//revive:enable:line-length-limit
}

func TestFlags(t *testing.T) {
	os.Args = []string{"cmd", "--dir=/home/directory", "--pkg=/go/pkg", "--out=raw_queries.go", "--fmt=false"}
	var cfg, err = loadCfg()

	if err != nil {
		t.Error(err)
	}

	if cfg.dir != "/home/directory" {
		t.Errorf("cfg.dir invalid: %s", cfg.dir)
	}

	if cfg.pkg != "/go/pkg" {
		t.Errorf("cfg.pkg invalid: %s", cfg.dir)
	}

	if cfg.out != "raw_queries.go" {
		t.Errorf("cfg.out invalid: %s", cfg.dir)
	}

	if cfg.fmt {
		t.Error("cfg.fmt invalid")
	}
}
