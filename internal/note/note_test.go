package note_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestNote(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("note", func(t *testing.T){
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testNoteContent", testNoteContent)
		t.Run("testNoteType", testNoteType)
		t.Run("testMemberNote", testMemberNote)
		t.Run("testNoteFirstAttachmentUrl", testNoteFirstAttachmentUrl)
	})
}

func testPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func testNoteContent(t *testing.T) {

	cases := []struct {
		ID     int
		Expect string
	}{
		{1, "Application note"},
		{2, "Issue raised"},
	}

	for _, c := range cases {
		r, err := note.ByID(db.Store, c.ID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.Expect, r.Content)
	}
}

func testNoteType(t *testing.T) {

	cases := []struct {
		ID     int
		Expect string
	}{
		{1, "General"},
		{2, "System"},
	}

	for _, c := range cases {
		r, err := note.ByID(db.Store, c.ID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.Expect, r.Type)
	}
}

func testMemberNote(t *testing.T) {

	xn, err := note.ByMemberID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(xn))
}

func testNoteFirstAttachmentUrl(t *testing.T) {

	cases := []struct {
		ID     int
		Expect string
	}{
		{1, "https://cdn.test.com/note/1/1-filename.ext"},
	}

	for _, c := range cases {
		n, err := note.ByID(db.Store, c.ID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.Expect, n.Attachments[0].URL)
	}
}
