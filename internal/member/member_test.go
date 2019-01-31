package member_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/testdata"
	"github.com/matryer/is"
	"gopkg.in/mgo.v2/bson"
)

var data = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestMember(t *testing.T) {

	err := data.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer data.TearDownMySQL()

	err = data.SetupMongoDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer data.TearDownMongoDB()

	t.Run("member", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testAddMember", testAddMember)
		t.Run("testByID", testByID)
		t.Run("testSearchDocDB", testSearchDocDB)
		t.Run("testSaveDocDB", testSaveDocDB)
		t.Run("testSyncUpdated", testSyncUpdated)
		t.Run("testExcelReport", testExcelReport)
		t.Run("testExcelReportJournal", testExcelReportJournal)
	})
}

func testPingDatabase(t *testing.T) {
	is := is.New(t)
	err := data.Store.MySQL.Session.Ping()
	is.NoErr(err) // Could not ping test database
}

// testAddMember tests the creation of a new member record
func testAddMember(t *testing.T) {
	m := member.Row{}
	m.RoleID = 2
	m.NamePrefixID = 1
	m.CountryID = 17
	m.ConsentDirectory = 1
	m.ConsentContact = 1
	m.UpdatedAt = "2019-01-01"
	m.DateOfBirth = "1970-11-03"
	m.Gender = "M"
	m.FirstName = "Mike"
	m.MiddleNames = "Peter"
	m.LastName = "Donnici"
	m.PostNominal = "B.Sc.Agr"
	m.QualificationsOther = "Grad. Cert. Computing"
	m.Mobile = "0402 400 191"
	m.PrimaryEmail = "michael@8o8.io"
	m.SecondaryEmail = "michael.donnici@gmail.com"

	m.QualificationRows = []member.QualificationRow{
		member.QualificationRow{
			MemberID:        m.ID,
			QualificationID: 11,
			OrganisationID:  222,
			YearObtained:    1992,
			Abbreviation:    "B.Sc.Agr.",
			Comment:         "Major in Crop Science",
		},
		member.QualificationRow{
			MemberID:        m.ID,
			QualificationID: 22,
			OrganisationID:  223,
			YearObtained:    1996,
			Abbreviation:    "Grad. Cert. Computing",
			Comment:         "Distance education",
		},
	}

	m.PositionRows = []member.PositionRow{
		member.PositionRow{
			MemberID:       m.ID,
			PositionID:     11,
			OrganisationID: 222,
			StartDate:      "2010-01-01",
			EndDate:        "2012-12-31",
			Comment:        "This is a comment",
		},
		member.PositionRow{
			MemberID:       m.ID,
			PositionID:     22,
			OrganisationID: 223,
			StartDate:      "2010-01-01",
			EndDate:        "2012-12-31",
			Comment:        "This is a comment",
		},
	}

	m.SpecialityRows = []member.SpecialityRow{
		member.SpecialityRow{
			MemberID:     m.ID,
			SpecialityID: 11,
			Preference:   1,
			Comment:      "This is a comment",
		},
	}

	err := m.Insert(data.Store)
	if err != nil {
		t.Fatalf("member.Row.Insert() err = %s", err)
	}
	if m.ID == 0 {
		t.Errorf("member.Row.ID = 0, want > 0")
	}

	// verify a few things about the member record
	mem, err := member.ByID(data.Store, m.ID)
	if err != nil {
		t.Fatalf("member.ByID(%d) err = %s", m.ID, err)
	}

	// check number of qualifications
	want := 2
	got := len(mem.Qualifications)
	if got != want {
		t.Errorf("member.Member.Qualifcations count = %d, want %d", got, want)
	}

	// check number of positions
	want = 2
	got = len(mem.Positions)
	if got != want {
		t.Errorf("member.Member.Positions count = %d, want %d", got, want)
	}

	// check number of specialities
	want = 1
	got = len(mem.Specialities)
	if got != want {
		t.Errorf("member.Member.Specialities count = %d, want %d", got, want)
	}
}

func testByID(t *testing.T) {
	is := is.New(t)
	m, err := member.ByID(data.Store, 1)
	is.NoErr(err)                                              // Error fetching member by id
	is.True(m.Active)                                          // Active should be true
	is.Equal(m.LastName, "Donnici")                            // Last name incorrect
	is.True(len(m.Memberships) > 0)                            // No memberships
	is.Equal(m.Memberships[0].Title, "Associate")              // Incorrect membership title
	is.Equal(m.Contact.EmailPrimary, "michael@mesa.net.au")    // Email incorrect
	is.Equal(m.Contact.Mobile, "0402123123")                   // Mobile incorrect
	is.Equal(m.Contact.Locations[0].City, "Jervis Bay")        // Location city incorrect
	is.Equal(m.Qualifications[0].Name, "PhD")                  // Qualification incorrect
	is.Equal(m.Specialities[1].Name, "Cardiac Cath Lab Nurse") // Speciality incorrect
	//printJSON(*m)
}

func testSearchDocDB(t *testing.T) {
	is := is.New(t)
	q := bson.M{"id": 7821}
	m, err := member.SearchDocDB(data.Store, q)
	is.NoErr(err)                     // Error querying MongoDB
	is.Equal(m[0].LastName, "Rousos") // Last name incorrect
}

func testSaveDocDB(t *testing.T) {
	is := is.New(t)
	mem := member.Member{
		ID:          1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Active:      true,
		Title:       "Mr",
		FirstName:   "Michael",
		MiddleNames: []string{"Peter"},
		LastName:    "Donnici",
		Gender:      "M",
		DateOfBirth: "1970-11-03",
	}
	err := mem.SaveDocDB(data.Store)
	is.NoErr(err) // Error saving to MongoDB

	q := bson.M{"lastName": "Donnici"}
	xm, err := member.SearchDocDB(data.Store, q)
	m := xm[0]
	is.NoErr(err)     // Error querying MongoDB
	is.Equal(m.ID, 1) // ID should be 1
}

func testSyncUpdated(t *testing.T) {
	is := is.New(t)
	mem := member.Member{
		ID:          2,
		CreatedAt:   time.Now().Add(-10 * time.Duration(time.Minute)), // 10 mins ago
		UpdatedAt:   time.Now().Add(-10 * time.Duration(time.Minute)), // 10 mins ago
		Active:      true,
		Title:       "Mr",
		FirstName:   "Barry",
		LastName:    "White",
		Gender:      "M",
		DateOfBirth: "1945-03-15",
	}
	err := mem.SaveDocDB(data.Store)
	is.NoErr(err) // Error saving to MongoDB

	memUpdate := member.Member{
		ID:          2,
		CreatedAt:   time.Now().Add(-10 * time.Duration(time.Minute)), // 10 mins ago
		UpdatedAt:   time.Now(),                                       // should trigger update
		Active:      false,
		Title:       "Mr",
		FirstName:   "Barry",
		LastName:    "White",
		Gender:      "M",
		DateOfBirth: "1948-03-15",
	}
	err = memUpdate.SyncUpdated(data.Store)
	is.NoErr(err) // Error syncing to MongoDB

	q := bson.M{"lastName": "White"}
	xm, err := member.SearchDocDB(data.Store, q)
	m := xm[0]
	is.NoErr(err)                         // Error querying MongoDB
	is.Equal(m.ID, 2)                     // ID should be 2
	is.Equal(m.Active, false)             // Active should be false
	is.Equal(m.DateOfBirth, "1948-03-15") // DateOfBirth incorrect
}

// fetch some test data and ensure excel report is not returning an error
func testExcelReport(t *testing.T) {

	id := 1   // member record
	want := 2 // expect 2 rows - heading and 2 record

	m, err := member.ByID(data.Store, id)
	if err != nil {
		t.Fatalf("member.ByID() err = %s", err)
	}
	xm := []member.Member{*m}
	f, err := member.ExcelReport(xm)
	if err != nil {
		t.Fatalf("member.ExcelReport() err = %s", err)
	}

	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}

// fetch some test data and ensure excel report (journal) is not returning an error
func testExcelReportJournal(t *testing.T) {

	id := 1   // member record
	want := 2 // expect 2 rows - heading and 2 record

	m, err := member.ByID(data.Store, id)
	if err != nil {
		t.Fatalf("member.ByID() err = %s", err)
	}
	xm := []member.Member{*m}
	f, err := member.ExcelReportJournal(xm)
	if err != nil {
		t.Fatalf("member.ExcelReportJournal() err = %s", err)
	}

	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}

func printJSON(m member.Member) {
	xb, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Print(string(xb))
	fmt.Println("-------------------------------------------------------------------")
}
