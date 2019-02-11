package speciality_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/speciality"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()

func TestSpeciality(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("Specialities", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testAll", testAll)
	})
}

func testPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("Ping() err = %s", err)
	}
}

// fetch the list of specialities
func testAll(t *testing.T) {
	xs, err := speciality.All(db.Store)
	if err != nil {
		t.Fatalf("speciality.All() err = %s", err)
	}
	got := len(xs)
	want := 5 // only 5 in test data
	if got != want {
		t.Errorf("speciality.All() count = %d, want %d", got, want)
	}
}
