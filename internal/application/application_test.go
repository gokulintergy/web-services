package application_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/application"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var store datastore.Datastore

func TestAll(t *testing.T) {

	var teardown func()
	store, teardown = setup()
	defer teardown()

	t.Run("application", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testNada", testNada)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("db.SetupMySQL() err = %s", err)
	}
	return db.Store, func() {
		err := db.TearDownMySQL()
		if err != nil {
			log.Fatalf("db.TearDownMySQL() err = %s", err)
		}
	}
}

func testPingDatabase(t *testing.T) {
	err := store.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("Ping() err = %s", err)
	}
}

func testNada(t *testing.T) {
	want := true
	got := application.Nada()
	if got != want {
		t.Errorf("Nada() = %v, want %v", got, want)
	}
}
