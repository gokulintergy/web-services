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

	// The following fields are values represented in junction tables
	QualificationRows []QualificationRow `json:"qualificationRows"`
	PositionRows      []PositionRow      `json:"positionRows"`
	SpecialityRows    []SpecialityRow    `json:"specialityRows"`
	AccreditationRows []AccreditationRow `json:"accreditationRows"`
	TagRows           []TagRow           `json:"tagRows"`

	// todo
	//Contact        Contact         `json:"contact"`
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

// PositionRow represents a member position in a junction table.
type PositionRow struct {
	ID             int    `json:"id"`
	MemberID       int    `json:"memberID"`
	PositionID     int    `json:"positionID"`
	OrganisationID int    `json:"organisationID"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	Comment        string `json:"comment"`
}

// SpecialityRow represents a member speciality in a junction table.
type SpecialityRow struct {
	ID           int    `json:"id"`
	MemberID     int    `json:"memberID"`
	SpecialityID int    `json:"specialityID"`
	Preference   int    `json:"preference"`
	Comment      string `json:"comment"`
}

// AccreditationRow represents a member accreditation in a junction table.
type AccreditationRow struct {
	ID              int    `json:"id"`
	MemberID        int    `json:"memberID"`
	AccreditationID int    `json:"accreditationID"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	Comment         string `json:"comment"`
}

// TagRow represents a member tag in a junction table.
type TagRow struct {
	ID       int `json:"id"`
	MemberID int `json:"memberID"`
	TagID    int `json:"tagID"`
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

	err = r.insertPositions(ds)
	if err != nil {
		return err
	}

	err = r.insertSpecialities(ds)
	if err != nil {
		return err
	}

	err = r.insertAccreditations(ds)
	if err != nil {
		return err
	}

	err = r.insertTags(ds)
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

// insertPositions inserts the member positions present in the Row value
func (r *Row) insertPositions(ds datastore.Datastore) error {
	for _, p := range r.PositionRows {
		err := p.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertSpecialities inserts the member specialities present in the Row value
func (r *Row) insertSpecialities(ds datastore.Datastore) error {
	for _, s := range r.SpecialityRows {
		err := s.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertAccreditations inserts the member accreditations present in the Row value
func (r *Row) insertAccreditations(ds datastore.Datastore) error {
	for _, a := range r.AccreditationRows {
		err := a.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertTags inserts the member tags present in the Row value
func (r *Row) insertTags(ds datastore.Datastore) error {
	for _, t := range r.TagRows {
		err := t.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insert a member qualification row in the junction table
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

// insert a member position row in the junction table
func (pr PositionRow) insert(ds datastore.Datastore, memberID int) error {
	query := fmt.Sprintf(queries["insert-member-position-row"],
		memberID,
		pr.PositionID,
		pr.OrganisationID,
		pr.StartDate,
		pr.EndDate,
		pr.Comment)
	_, err := ds.MySQL.Session.Exec(query)
	return err
}

// insert a member speciality row in the junction table
func (sr SpecialityRow) insert(ds datastore.Datastore, memberID int) error {
	query := fmt.Sprintf(queries["insert-member-speciality-row"],
		memberID,
		sr.SpecialityID,
		sr.Preference,
		sr.Comment)
	_, err := ds.MySQL.Session.Exec(query)
	return err
}

// insert a member accreditation row in the junction table
func (ar AccreditationRow) insert(ds datastore.Datastore, memberID int) error {
	query := fmt.Sprintf(queries["insert-member-accreditation-row"],
		memberID,
		ar.AccreditationID,
		ar.StartDate,
		ar.EndDate,
		ar.Comment)
	_, err := ds.MySQL.Session.Exec(query)
	return err
}

// insert a member tag row in the junction table
func (tr TagRow) insert(ds datastore.Datastore, memberID int) error {
	query := fmt.Sprintf(queries["insert-member-tag-row"],
		memberID,
		tr.TagID)
	_, err := ds.MySQL.Session.Exec(query)
	return err
}
