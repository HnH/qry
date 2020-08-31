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

	if len(q) != 3 {
		t.Error("Expected 3 files")
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

	if _, ok := q["three.sql"]; !ok {
		t.Error("three.sql not loaded")
	}

	if q["three.sql"]["EscapedJSONQuery"] != "INSERT INTO \\\"data\\\" (id, \\\"data\\\") VALUES (1, '{\\\"test\\\": 1}'), (2, '{\\\"test\\\": 2}');" {
		t.Error("Invalid EscapedJSONQuery query")
	}

	if q["three.sql"]["EscapedByteaQuery"] != "INSERT INTO bin (id, \\\"data\\\") VALUES (1, E'\\\\\\\\x3aaab6e7fb7245cc9785653a0d9ffc4a5ce0f974');" {
		t.Error("Invalid EscapedByteaQuery query")
	}
}

func TestDirInvalid(t *testing.T) {
	var _, err = Dir("11")
	if err == nil {
		t.Error("Error expected on invalid directory")
	}
}
