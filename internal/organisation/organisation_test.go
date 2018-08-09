package organisation_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/cardiacsociety/web-services/internal/organisation"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestOrganisation(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("organisation", func(t *testing.T){
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testOrganisationByID", testOrganisationByID)
		t.Run("testOrganisationDeepEqual", testOrganisationDeepEqual)
		t.Run("testOrganisationCount", testOrganisationCount)
		t.Run("testChildOrganisationCount", testChildOrganisationCount)
	})
}

func testPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func testOrganisationByID(t *testing.T) {
	org, err := organisation.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, "ABC Organisation", org.Name)
}

func testOrganisationDeepEqual(t *testing.T) {

	exp := organisation.Organisation{
		ID:   1,
		Code: "ABC",
		Name: "ABC Organisation",
		Groups: []organisation.Organisation{
			{ID: 3, Code: "ABC-1", Name: "ABC Sub1"},
			{ID: 4, Code: "ABC-2", Name: "ABC Sub2"},
			{ID: 5, Code: "ABC-3", Name: "ABC Sub3"},
		},
	}

	o, err := organisation.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	res := reflect.DeepEqual(exp, o)
	helper.Result(t, true, res)
}

// Test data has 2 parent organisations
func testOrganisationCount(t *testing.T) {
	xo, err := organisation.All(db.Store)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 2, len(xo))
}

// Test data has 3 child organisations belonging to parent id 1
func testChildOrganisationCount(t *testing.T) {
	o, err := organisation.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(o.Groups))
}
