package db_test

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	"io/ioutil"
	"os"
)

func newTestLDB() (*db.LDBDatabase, func()) {
	dirname, err := ioutil.TempDir(os.TempDir(), "db_test_")
	if err != nil {
		panic("failed to create test file: " + err.Error())
	}
	db, err := db.NewLDBDatabase(dirname, 0, 0)
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}
	return db, func() {
		db.Close()
		os.RemoveAll(dirname)
	}
}