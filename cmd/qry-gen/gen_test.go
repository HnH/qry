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
			dir: "../../testdata",
			out: "gen.go",
			fmt: true,
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

	for i, c := range []bool{
		bytes.Contains(b, []byte("package pkg")),
		bytes.Contains(b, []byte("\t// one.sql")),
		bytes.Contains(b, []byte("\tInsertUser  = \"INSERT INTO `users` (`name`) VALUES (?);\"")),
		bytes.Contains(b, []byte("\tGetUserById = \"SELECT * FROM `users` WHERE `user_id` = ?;\"")),
		bytes.Contains(b, []byte("\t// two.sql")),
		bytes.Contains(b, []byte("\tDeleteUsersByIds = \"DELETE FROM `users` WHERE `user_id` IN ({ids});\"")),
	} {
		if !c {
			t.Errorf("check at idx %d has not passed", i)
		}
	}
}
