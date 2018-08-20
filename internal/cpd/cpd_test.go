package cpd_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestCPD(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("CPD", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testCPDByID", testCPDByID)
		t.Run("testCPDByMemberID", testCPDByMemberID)
		t.Run("testCPDQuery", testCPDQuery)
		t.Run("testAddCPD", testAddCPD)
		t.Run("testUpdateCPD", testUpdateCPD)
		t.Run("testDuplicateOf", testDuplicateOf)
		t.Run("testDelete", testDelete)
	})
}

func testPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func testCPDByID(t *testing.T) {

	cases := []struct {
		id   int
		desc string
	}{
		{1, "BJJ like Bruno Malfacine"},
		{2, "Ate sausages and eggs"},
		{3, "Baked bread"},
	}

	for _, c := range cases {
		cpd, err := cpd.ByID(db.Store, c.id)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.desc, cpd.Description)
	}
}

func testCPDByMemberID(t *testing.T) {
	xcpd, err := cpd.ByMemberID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(xcpd))
}

func testCPDQuery(t *testing.T) {
	xcpd, err := cpd.Query(db.Store, "WHERE cma.description LIKE '%Bruno%'")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, len(xcpd))
}

func testAddCPD(t *testing.T) {
	c := cpd.Input{
		MemberID:    1,
		ActivityID:  24,
		TypeID:      25,
		Date:        "2018-05-07",
		Quantity:    2.25,
		Description: "I added this record",
		Evidence:    false,
	}
	id, err := cpd.Add(db.Store, c)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	// fetch the newly added record, and verify the description
	r, err := cpd.ByID(db.Store, id)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	helper.Result(t, c.Description, r.Description)
}

func testUpdateCPD(t *testing.T) {
	c := cpd.Input{
		ID:          2,
		MemberID:    1,
		ActivityID:  24,
		TypeID:      25,
		Date:        "2018-05-07",
		Quantity:    2.25,
		Description: "The description was updated",
		Evidence:    false,
	}
	err := cpd.Update(db.Store, c)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	r, err := cpd.ByID(db.Store, c.ID)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	helper.Result(t, c.Description, r.Description)
}

func testDuplicateOf(t *testing.T) {

	// Fetch first cpd record and then try to insert it - should get '1' returned
	a, err := cpd.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	i := cpd.Input{
		MemberID:    a.MemberID,
		ActivityID:  a.Activity.ID,
		TypeID:      int(a.Type.ID.Int64),
		Date:        a.Date,
		Description: a.Description,
		Evidence:    a.Evidence,
		UnitCredit:  a.CreditData.UnitCredit,
		Quantity:    a.Credit,
	}

	dupID, err := cpd.DuplicateOf(db.Store, i)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	helper.Result(t, 1, dupID)
}

func testDelete(t *testing.T) {

	// get count, delete record id 3, count should be count - 1
	xcpd, err := cpd.Query(db.Store, "")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	c1 := len(xcpd)

	err = cpd.Delete(db.Store, 1, 3)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	xcpd, err = cpd.Query(db.Store, "")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	c2 := len(xcpd)

	helper.Result(t, c1-1, c2)
}
