package issue_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/issue"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestIssue(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("issue", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testInsertRowErrorIDNotNil", testInsertRowErrorIDNotNil)
		t.Run("testInsertRowErrorNoTypeID", testInsertRowErrorNoTypeID)
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
	err := ds.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("Ping() err = %s", err)
	}
}

// test an attempt to insert an issue row when the Issue.ID has a value
func testInsertRowErrorIDNotNil(t *testing.T) {
	i := issue.Issue{
		ID:          1,
		Type:        issue.Type{ID: 2},
		Description: "This is the description",
		Action:      "This is what must be done",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorIDNotNil
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %s", want)
	}
}

// test an attempt to insert an issue row with Issue.Type.ID not set
func testInsertRowErrorNoTypeID(t *testing.T) {
	i := issue.Issue{
		Description: "This is the description",
		Action:      "This is what must be done",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorNoTypeID
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %s", want)
	}
}
