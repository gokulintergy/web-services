package member

import (
	"encoding/json"
	"strings"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Row represents a raw record from the member table in the SQL database. This type is
// primarily for inserting new records. Junction table data are represented with []int containing
// a list of foreign key ids for the relevant table. The comment next to JSON tags match the
// relational db columns names.
type Row struct {
	ID                 int    `json:"id"`
	RoleID             int    `json:"roleId"`             // acl_member_role_id
	NamePrefixID       int    `json:"titleId"`            // a_name_prefix_id
	CountryID          int    `json:"countryId"`          // country_id
	ConsentDirectory   bool   `json:"consentDirectory"`   // consent_directory
	ConsentContact     bool   `json:"consentContact"`     // consent_contact
	UpdatedAt          string `json:"updatedAt"`          // updated_at
	DateOfBirth        string `json:"dateOfBirth"`        // date_of_birth
	Gender             string `json:"gender"`             // gender
	FirstName          string `json:"firstName"`          // first_name
	MiddleNames        string `json:"middleNames"`        // middle_names
	LastName           string `json:"lastName"`           // last_name
	PostNominal        string `json:"postNominal"`        // suffix
	QualificationsInfo string `json:"qualificationsInfo"` // qualifications_other
	Mobile             string `json:"mobile"`             // mobile_phone
	PrimaryEmail       string `json:"primaryEmail"`       // primary_email
	SecondaryEmail     string `json:"secondaryEmail"`     // secondary_email

	// The following fields are values represented in junction tables
	QualificationRows []QualificationRow `json:"qualifications"`
	SpecialityRows    []SpecialityRow    `json:"interests"`
	PositionRows      []PositionRow      `json:"positions"`
	AccreditationRows []AccreditationRow `json:"accreditations"`
	TagRows           []TagRow           `json:"tags"`
	ApplicationRow    ApplicationRow     `json:"application"`
	ContactRows       []ContactRow       `json:"contacts"`
}

// QualificationRow represents a member qualification in a junction table.
type QualificationRow struct {
	ID              int    `json:"id"` // id of the junction record
	MemberID        int    `json:"memberId"`
	QualificationID int    `json:"qualificationId"`
	OrganisationID  int    `json:"organisationId"`
	YearObtained    int    `json:"year"`
	Abbreviation    string `json:"abbreviation"`
	Comment         string `json:"comment"`
}

// PositionRow represents a member position in a junction table.
type PositionRow struct {
	ID             int    `json:"id"` // id of the junction record
	MemberID       int    `json:"memberId"`
	PositionID     int    `json:"positionId"`
	OrganisationID int    `json:"organisationId"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	Comment        string `json:"comment"`
}

// SpecialityRow represents a member speciality in a junction table.
type SpecialityRow struct {
	ID           int    `json:"id"`
	MemberID     int    `json:"memberId"`
	SpecialityID int    `json:"specialityId"`
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

// ApplicationRow represents a member's application
type ApplicationRow struct {
	ID          int
	MemberID    int
	ForTitleID  int    `json:"forTitleId"`
	NominatorID int    `json:"nominatorId"`
	SeconderID  int    `json:"seconderId"`
	Comment     string `json:"nominatorInfo"`
}

// ContactRow represents a contact location
type ContactRow struct {
	ID        int
	MemberID  int
	TypeID    int    `json:"contactTypeId"`
	Phone     string `json:"phone"`
	Fax       string `json:"fax"`
	Email     string `json:"email"`
	Web       string `json:"web"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Address3  string `json:"address3"`
	Locality  string `json:"locality"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	CountryID int    `json:"countryId"`
}

// Insert inserts a member row into the database. If successful it will set the member id.
func (r *Row) Insert(ds datastore.Datastore) error {

	// convert bools to 0/1
	var consentDirectory, consentContact int
	if r.ConsentDirectory {
		consentDirectory = 1
	}
	if r.ConsentContact {
		consentContact = 1
	}

	// gender stored as 'M' or 'F', so capitalise first letter of gender string
	r.Gender = strings.ToUpper(string(strings.TrimSpace(r.Gender)[0]))

	res, err := ds.MySQL.Session.Exec(queries["insert-member-row"],
		r.RoleID,
		r.NamePrefixID,
		r.CountryID,
		consentDirectory,
		consentContact,
		r.DateOfBirth,
		r.Gender,
		r.FirstName,
		r.MiddleNames,
		r.LastName,
		r.PostNominal,
		r.QualificationsInfo,
		r.Mobile,
		r.PrimaryEmail,
		r.SecondaryEmail)
	if err != nil {
		return err
	}

	// get member id
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

	err = r.insertApplication(ds)
	if err != nil {
		return err
	}

	err = r.insertContacts(ds)
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

// insertApplication creates an application record for the member
func (r *Row) insertApplication(ds datastore.Datastore) error {
	return r.ApplicationRow.insert(ds, r.ID)
}

// insertContacts inserts the member contact rows
func (r *Row) insertContacts(ds datastore.Datastore) error {
	for _, c := range r.ContactRows {
		err := c.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insert a member qualification row in the junction table
func (qr QualificationRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-qualification-row"],
		memberID,
		qr.QualificationID,
		qr.OrganisationID,
		qr.YearObtained,
		qr.Abbreviation,
		qr.Comment)
	return err
}

// insert a member position row in the junction table
func (pr PositionRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-position-row"],
		memberID,
		pr.PositionID,
		pr.OrganisationID,
	)
	return err
}

// insert a member speciality row in the junction table
func (sr SpecialityRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-speciality-row"],
		memberID,
		sr.SpecialityID,
		sr.Preference,
		sr.Comment)
	return err
}

// insert a member accreditation row in the junction table
func (ar AccreditationRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-accreditation-row"],
		memberID,
		ar.AccreditationID,
		ar.StartDate,
		ar.EndDate,
		ar.Comment)
	return err
}

// insert a member tag row in the junction table
func (tr TagRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-tag-row"],
		memberID,
		tr.TagID)
	return err
}

func (ar ApplicationRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-application-row"],
		memberID,
		ar.NominatorID,
		ar.SeconderID,
		ar.ForTitleID,
		ar.Comment)
	return err
}

func (cr ContactRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-contact-row"],
		memberID,
		cr.TypeID,
		cr.CountryID,
		cr.Phone,
		cr.Fax,
		cr.Email,
		cr.Web,
		cr.Address1,
		cr.Address2,
		cr.Address3,
		cr.Locality,
		cr.State,
		cr.Postcode,
	)
	return err
}

// InsertRowFromJSON creates a member Row from a JSON object. This is used for online
// membership applications and returns a new member Row on success.
func InsertRowFromJSON(ds datastore.Datastore, s string) (Row, error) {
	r := Row{}

	err := json.Unmarshal([]byte(s), &r)
	if err != nil {
		return r, err
	}
	err = r.Insert(ds)
	if err != nil {
		return r, err
	}
	return r, nil
}
