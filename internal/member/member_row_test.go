package member_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

// Note names: ds2 and setup2() as package-level identifiers must be unique -
// ds and setup() exist in member_test.go
var ds2 datastore.Datastore

func TestMemberRow(t *testing.T) {

	//var teardown func()
	ds2, _ = setup2()
	//defer teardown()

	t.Run("member_row", func(t *testing.T) {
		//t.Run("testInsertRow", testInsertRow)
		t.Run("testInsertRowJSON", testInsertRowJSON)
	})
}

func setup2() (datastore.Datastore, func()) {
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

// testInsertRow tests the creation of a new member record
func testInsertRow(t *testing.T) {
	m := member.Row{}
	m.RoleID = 2
	m.NamePrefixID = 1
	m.CountryID = 17
	m.ConsentDirectory = true
	m.ConsentContact = true
	m.UpdatedAt = "2019-01-01"
	m.DateOfBirth = "1970-11-03"
	m.Gender = "M"
	m.FirstName = "Mike"
	m.MiddleNames = "Peter"
	m.LastName = "Donnici"
	m.PostNominal = "B.Sc.Agr"
	m.QualificationsInfo = "Grad. Cert. Computing"
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

	m.AccreditationRows = []member.AccreditationRow{
		member.AccreditationRow{
			MemberID:        m.ID,
			AccreditationID: 11,
			StartDate:       "2010-01-01",
			EndDate:         "2012-12-31",
			Comment:         "This is a comment",
		},
	}

	m.TagRows = []member.TagRow{
		member.TagRow{
			MemberID: m.ID,
			TagID:    1,
		},
		member.TagRow{
			MemberID: m.ID,
			TagID:    2,
		},
		member.TagRow{
			MemberID: m.ID,
			TagID:    3,
		},
	}

	m.ContactRows = []member.ContactRow{
		member.ContactRow{
			MemberID: m.ID,
			TypeID: 2, // Directory
			CountryID: 14,
			Phone: "02 444 66 789",
			Fax: "02 444 66 890",
			Email: "any@oldemail.com",
			Web: "https://thesite.com",
			Address1: "Leve 12",
			Address2: "123 Some Street",
			Address3: "Some large building",
			Locality: "CityTown",
			State: "NewShire",
			Postcode: "1234",
		},
		member.ContactRow{
			MemberID: m.ID,
			TypeID: 1, // Mail
			CountryID: 14,
			Address1: "Level 12",
			Address2: "123 Some Street",
			Address3: "Some large building",
			Locality: "CityTown",
			State: "NewShire",
			Postcode: "1234",
		},
	}	

	err := m.Insert(ds2)
	if err != nil {
		t.Fatalf("member.Row.Insert() err = %s", err)
	}
	if m.ID == 0 {
		t.Errorf("member.Row.ID = 0, want > 0")
	}

	// verify a few things about the member record
	mem, err := member.ByID(ds2, m.ID)
	if err != nil {
		t.Fatalf("member.ByID(%d) err = %s", m.ID, err)
	}

	// check number of qualifications
	want := 2
	got := len(mem.Qualifications)
	if got != want {
		t.Errorf("Member.Qualifcations count = %d, want %d", got, want)
	}

	// check number of positions
	want = 2
	got = len(mem.Positions)
	if got != want {
		t.Errorf("Member.Positions count = %d, want %d", got, want)
	}

	// check number of specialities
	want = 1
	got = len(mem.Specialities)
	if got != want {
		t.Errorf("Member.Specialities count = %d, want %d", got, want)
	}

	// check number of accreditations
	want = 1
	got = len(mem.Accreditations)
	if got != want {
		t.Errorf("Member.Accreditations count = %d, want %d", got, want)
	}

	// check number of tags
	want = 3
	got = len(mem.Tags)
	if got != want {
		t.Errorf("Member.Tags count = %d, want %d", got, want)
	}

	// check number of contacts
	want = 2
	got = len(mem.Contact.Locations)
	if got != want {
		t.Errorf("Member.Contact.Locations count = %d, want %d", got, want)
	}
}

// testInsertRowJSON tests the creation of a new member record from a JSON doc
func testInsertRowJSON(t *testing.T) {

	// When this test is passing, below is the format for JSON posted to create a new application
	j := `{
		"roleId" : 2,
		"countryId": 14, 
		"gender": "Male",
		"title": "Dr",
		"titleId": 5,
		"firstName": "Mike",
		"middleNames": "Peter",
		"lastName": "Donnici",
		"dateOfBirth": "1970-11-03",
		"primaryEmail": "michael@somewhere.com",
		"secondaryEmail": "michael@somewhereelse.com",
		"mobile": "+61402400191",
		"consentDirectory": true,
		"consentContact": true,

		"contacts": [
			{
				"contactTypeId": 2,
				"countryId": 14,
				"phone": "02 444 66 789",
				"fax": "02 444 66 890",
				"email": "any@oldemail.com",
				"web": "https://thesite.com",
				"address1": "Level 12",
				"address2": "123 Some Street",
				"address3": "Some large building",
				"locality": "CityTown",	
				"state": "NewShire",
				"postcode": "1234"
			},
			{
				"contactTypeId": 1,
				"countryId": 14,
				"phone": "02 444 66 789",
				"fax": "02 444 66 890",
				"address1": "Level 12",
				"address2": "123 Some Street",
				"address3": "Some large building",
				"locality": "CityTown",	
				"state": "NewShire",
				"postcode": "1234"
			}
		],

		"trainee": true,
		"tags" : [
			{"tagId": 4}
		],

		"qualifications": [
			{
				"qualificationId": 2,
				"name": "Bachelor of Medicine, Bachelor of Surgery",
				"abbreviation": "MBBS",
				"year": 1998,
				"organisationId": 311,
				"organisationName": "University of Adelaide"
			}
		],
		"qualificationsInfo": "ABC123",

		"interests": [
			{
				"specialityId": 1,
				"name": "Cardiac Care Nurse (Medical)"
			}, 
			{
				"specialityId": 2,
				"name": "Cardiac Cath Lab Nurse"
			}
		],

		"positions": [
			{
				"positionId": 1,
				"organisationId": 4
			},
			{
				"positionId": 2,
				"organisationId": 2
			},
			{
				"positionId": 3,
				"organisationId": 11
			} 
		],

		"application": {
			"forTitle": "Associate",
			"forTitleId": 2,
			"nominatorId": 586,
			"seconderId": 587,
			"nominatorInfo": "ghggh"
		}, 
		
		"ishr": false,
		"consentRequestInfo": true
	}`

	row, err := member.InsertRowFromJSON(ds2, j)
	if err != nil {
		t.Fatalf("member.RowFromJSON() err = %s", err)
	}

	// verify a few things about the member record
	mem, err := member.ByID(ds2, row.ID)
	if err != nil {
		t.Fatalf("member.ByID(%d) err = %s", row.ID, err)
	}

	// check number of qualifications
	want := 1
	got := len(mem.Qualifications)
	if got != want {
		t.Errorf("Member.Qualifications count = %d, want %d", got, want)
	}

	wantQualOther := "ABC123"
	gotQualOther := mem.QualificationsOther
	if gotQualOther != wantQualOther {
		t.Errorf("Member.QualificationsOther = %q, want %q", gotQualOther, wantQualOther)
	}

	// check number of positions
	want = 3
	got = len(mem.Positions)
	if got != want {
		t.Errorf("Member.Positions count = %d, want %d", got, want)
	}

	// check number of specialities
	want = 2
	got = len(mem.Specialities)
	if got != want {
		t.Errorf("Member.Specialities count = %d, want %d", got, want)
	}

	// check number of tags
	want = 1
	got = len(mem.Tags)
	if got != want {
		t.Errorf("Member.Tags count = %d, want %d", got, want)
	}

	// check number of contacts
	want = 2
	got = len(mem.Contact.Locations)
	if got != want {
		t.Errorf("Member.Contact.Locations count = %d, want %d", got, want)
	}
}
