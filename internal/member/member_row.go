package member

import (
	"fmt"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Row represents a raw record from the member table in the SQL database. This type is
// primarily for inserting new records. Junction table data are represented with []int containing
// a list of foreign key ids for the relevant table. The JSON tags match the columns names.
type Row struct {
	ID                  int    `json:"id"`
	RoleID              int    `json:"acl_member_role_id"`
	NamePrefixID        int    `json:"a_name_prefix_id"`
	CountryID           int    `json:"country_id"`
	ConsentDirectory    int    `json:"consent_directory"`
	ConsentContact      int    `json:"consent_contact"`
	UpdatedAt           string `json:"updatedAt"`
	DateOfBirth         string `json:"date_of_birth"`
	Gender              string `json:"gender"`
	FirstName           string `json:"first_name"`
	MiddleNames         string `json:"middle_names"`
	LastName            string `json:"last_name"`
	PostNominal         string `json:"suffix"`
	QualificationsOther string `json:"qualifications_other"`
	Mobile              string `json:"mobile_phone"`
	PrimaryEmail        string `json:"primary_email"`
	SecondaryEmail      string `json:"secondary_email"`

	// The following fields are values repreented in junction tables
	//Contact        Contact         `json:"contact"`
	QualificationRows []QualificationRow `json:"qualificationRows"`
	AccreditationIDs  []int              `json:"accreditations"`
	PositionIDs       []int              `json:"positions"`
	SpecialityIDs     []int              `json:"specialities"`
	TagIDs            []int              `json:"tags"`
}

// QualificationRow represents a member qualification in a junction table.
type QualificationRow struct {
	ID              int    `json:"id"`
	MemberID        int    `json:"memberID"`
	QualificationID int    `json:"qualificationID"`
	OrganisationID  int    `json:"organisationID"`
	YearObtained    int    `json:"yearObtained"`
	Abbreviation    string `json:"abbreviation"`
	Comment         string `json:"comment"`
}

// Insert inserts a member row into the database. If successful it will set the member id.
func (r *Row) Insert(ds datastore.Datastore) error {
	query := fmt.Sprintf(queries["insert-member-row"],
		r.RoleID,
		r.NamePrefixID,
		r.CountryID,
		r.ConsentDirectory,
		r.ConsentContact,
		r.DateOfBirth,
		r.Gender,
		r.FirstName,
		r.MiddleNames,
		r.LastName,
		r.PostNominal,
		r.QualificationsOther,
		r.Mobile,
		r.PrimaryEmail,
		r.SecondaryEmail,
	)
	res, err := ds.MySQL.Session.Exec(query)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = int(id) // from int64

	err = r.insertQualifications(ds)
	if err != nil {
		return err
	}

	return nil
}

// insertQualifications inserts the member qualifications present in the Row value
func (r *Row) insertQualifications(ds datastore.Datastore) error {
	for _, q := range r.QualificationRows {
		err := q.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertMemberQualification inserts a member qualification via a junction table
func (qr QualificationRow) insert(ds datastore.Datastore, memberID int) error {
	query := fmt.Sprintf(queries["insert-member-qualification-row"],
		memberID,
		qr.QualificationID,
		qr.OrganisationID,
		qr.YearObtained,
		qr.Abbreviation,
		qr.Comment)
	_, err := ds.MySQL.Session.Exec(query)
	return err
}
