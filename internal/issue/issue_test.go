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

	// var teardown func()
	ds, _ = setup()
	// defer teardown()

	t.Run("issue", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testInsertRowErrorIDNotNil", testInsertRowErrorIDNotNil)
		t.Run("testInsertRowErrorNoTypeID", testInsertRowErrorNoTypeID)
		t.Run("testInsertRowErrorNoDescription", testInsertRowErrorNoDescription)
		t.Run("testInsertRowErrorAssociationNoMemberID", testInsertRowErrorAssociationNoMemberID)
		t.Run("testInsertRowErrorAssociation", testInsertRowErrorAssociation)
		t.Run("testInsertRowErrorAssociationID", testInsertRowErrorAssociationID)
		t.Run("testInsertRowErrorAssociationEntity", testInsertRowErrorAssociationEntity)
		t.Run("testInsertRow", testInsertRow)
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
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
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
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

// test an attempt to insert an issue row with Issue.Description not set
func testInsertRowErrorNoDescription(t *testing.T) {
	i := issue.Issue{
		Type:   issue.Type{ID: 2},
		Action: "This is what must be done",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorNoDescription
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociationNoMemberID(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      0, // err
		AssociationID: 345,
		Association:   "application",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociationNoMemberID
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociation(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 345,
		Association:   "", // err
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociation
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociationID(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 0, // err
		Association:   "application",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociationID
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociationEntity(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 345,
		Association:   "unknownentity", // err
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociationEntity
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

// test insert an error row
func testInsertRow(t *testing.T) {
	i := issue.Issue{
		Type:        issue.Type{ID: 2},
		Description: "This is the description",
		Action:      "This is what must be done",
	}
	err := i.InsertRow(ds)
	if err != nil {
		t.Errorf("Issue.InsertRow() err = %s", err)
	}
}
