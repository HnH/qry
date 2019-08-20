package qry

import (
	"testing"
)

func TestQuery(t *testing.T) {
	var expected = Query("DELETE FROM `users` WHERE `user_id` IN (?,?,?);")
	if Query("DELETE FROM `users` WHERE `user_id` IN ({ids});").Replace("{ids}", In(3)) != expected {
		t.Error("Invalid qry.In() result")
	}
}

func TestReplaceEmpty(t *testing.T) {
	var expected = Query("DELETE FROM `users` WHERE `user_id` IN ({ids});")
	if expected.Replace("", "") != expected {
		t.Error("Expected empty qry.Replace() result")
	}
}

func TestInEmpty(t *testing.T) {
	if In(-1) != "" {
		t.Error("Expected empty qry.In() result")
	}
}
