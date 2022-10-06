package database

import "testing"

func TestDB(t *testing.T) {
	db := NewDB()
	db.Insert("foo", "bar")
	db.Insert("this", "===")
	db.Insert("baz", "doo")
}
