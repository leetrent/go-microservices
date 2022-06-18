package main

import (
	"authentication/data"
	"os"
	"testing"
)

// func TestMain(m *testing.M} {
// 	os.Exit(m.Run())
// }

var testApp Config

func TestMain(m *testing.M) {
	repo := data.NewPostgresTestRepository(nil)
	testApp.Repo = repo
	os.Exit(m.Run())
}
