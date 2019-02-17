package note_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestNote(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("note", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testNoteContent", testNoteContent)
		t.Run("testNoteType", testNoteType)
		t.Run("testMemberNote", testMemberNote)
		t.Run("testNoteFirstAttachmentUrl", testNoteFirstAttachmentUrl)
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
		t.Fatal("Could not ping database")
	}
}

func testNoteContent(t *testing.T) {
	cases := []struct {
		arg  int
		want string
	}{
		{1, "Application note"},
		{2, "Issue raised"},
	}
	for _, c := range cases {
		n, err := note.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("note.ByID(%d) err = %s", c.arg, err)
		}
		got := n.Content
		if got != c.want {
			t.Errorf("Note.Content = %q, want %q", got, c.want)
		}
	}
}

func testNoteType(t *testing.T) {
	cases := []struct {
		arg  int
		want string
	}{
		{1, "General"},
		{2, "System"},
	}
	for _, c := range cases {
		n, err := note.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("note.ByID(%d) err = %s", c.arg, err)
		}
		got := n.Type
		if got != c.want {
			t.Errorf("Note.Type = %q, want %q", got, c.want)
		}
	}
}

func testMemberNote(t *testing.T) {
	arg := 1 // member id
	xn, err := note.ByMemberID(ds, 1)
	if err != nil {
		t.Fatalf("note.ByMemberID(%d) err = %s", arg, err)
	}
	got := len(xn)
	want := 3
	if got != want {
		t.Errorf("note.ByMemberID(%d) count = %d, want %d", arg, got, want)
	}
}

func testNoteFirstAttachmentUrl(t *testing.T) {
	cases := []struct {
		arg  int    // note id
		want string // file url
	}{
		{1, "https://cdn.test.com/note/1/1-filename.ext"},
	}

	for _, c := range cases {
		n, err := note.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("note.ByID(%d) err = %s", c.arg, err)
		}
		got := n.Attachments[0].URL
		if got != c.want {
			t.Errorf("Note.Attachments[0].URL = %s, want %s", got, c.want)
		}
	}
}
