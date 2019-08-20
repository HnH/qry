package qry

import (
	"testing"
)

func TestDir(t *testing.T) {
	var (
		q   map[string]QuerySet
		err error
	)

	if q, err = Dir("./testdata"); err != nil {
		t.Error(err)
	}

	if len(q) != 2 {
		t.Error("Expected 2 files")
	}

	if _, ok := q["one.sql"]; !ok {
		t.Error("one.sql not loaded")
	}

	if q["one.sql"]["InsertUser"] != "INSERT INTO `users` (`name`) VALUES (?);" {
		t.Error("Invalid InsertUser query")
	}

	if q["one.sql"]["GetUserById"] != "SELECT * FROM `users` WHERE `user_id` = ?;" {
		t.Error("Invalid GetUserById query")
	}

	if _, ok := q["two.sql"]; !ok {
		t.Error("two.sql not loaded")
	}

	if q["two.sql"]["DeleteUsersByIds"] != "DELETE FROM `users` WHERE `user_id` IN ({ids});" {
		t.Error("Invalid DeleteUsersByIds query")
	}
}

func TestDirInvalid(t *testing.T) {
	var _, err = Dir("11")
	if err == nil {
		t.Error("Error expected on invalid directory")
	}
}