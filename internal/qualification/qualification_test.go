package qualification_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/qualification"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestQualification(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("Qualifications", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testAll", testAll)
	})
}

func testPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

// fetch the list of qualifications
func testAll(t *testing.T) {
	xq, err := qualification.All(db.Store)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 29, len(xq))
}
