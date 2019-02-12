package tag_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/tag"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()

func TestTag(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("Tags", func(t *testing.T) {
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

// fetch the list of tags
func testAll(t *testing.T) {
	xs, err := tag.All(db.Store)
	if err != nil {
		t.Fatalf("Tag.All() err = %s", err)
	}
	got := len(xs)
	want := 5 // only 5 in test data
	if got != want {
		t.Errorf("Tag.All() count = %d, want %d", got, want)
	}
}
