// Package note provides management of notes data
package note

import (
	"errors"
	"fmt"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Error messages
const (
	ErrorIDNotNil              = "cannot insert a note row because Note.ID already has a value"
	ErrorNoMemberID            = "cannot insert a note row because Note.MemberID field is not set"
	ErrorNoTypeID              = "cannot insert a note row because Note.TypeID field is not set"
	ErrorNoContent             = "cannot insert a note row because Note.Content is empty"
	ErrorAssociation           = "association entity not specified"
	ErrorAssociationID         = "association entity ID not specified"
	ErrorAssociationEntity     = "association entity invalid"
)

// Note represents a record of a comment, document or anything else. A Note is always linked to a member 
// and can also be associated with an application or an issue
type Note struct {
	ID            int `json:"id" bson:"id"`
	MemberID      int          `json:"memberId" bson:"memberId"`
	TypeID        int
	Type          string       `json:"type" bson:"type"`
	Association   string // either "application" or "invoice"
	AssociationID int    // the id of the associated application or invoice record
	DateCreated   string       `json:"dateCreated" bson:"dateCreated"`
	DateUpdated   string       `json:"dateUpdated" bson:"dateUpdated"`
	DateEffective string       `json:"dateEffective" bson:"dateEffective"`
	Content       string       `json:"content" bson:"content"`
	Attachments   []Attachment `json:"attachments" bson:"attachments"`
}

// Attachment is a file linked to a note
type Attachment struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// InsertRow creates a new note row with fields from Note
func (n *Note) InsertRow(ds datastore.Datastore) error {
	switch {
	case n.ID > 0:
		return errors.New(ErrorIDNotNil)
	case n.MemberID == 0:
		return errors.New(ErrorNoMemberID)
	case n.TypeID == 0:
		return errors.New(ErrorNoTypeID)
	case n.Content == "":
		return errors.New(ErrorNoContent)
	}
	q := fmt.Sprintf(queries["insert-note"], n.TypeID, n.Content)
	res, err := ds.MySQL.Session.Exec(q)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = int(id) // from int64

	// Associate other data if fields are set
	if n.AssociationID > 0 || n.Association != "" {
		err := n.checkAssociatioData()
		if err != nil {
			return err
		}
		q := fmt.Sprintf(queries["insert-note-association"], n.ID, n.MemberID, n.AssociationID, n.Association)
		_, err = ds.MySQL.Session.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkAssociatioData verifies fields required to associate an issue with other data
func (n *Note) checkAssociatioData() error {
	// no association
	if n.Association == "" && n.AssociationID == 0 {
		return nil
	}
	if n.Association != "" || n.AssociationID > 0 {
		switch {
		case n.Association == "":
			return errors.New(ErrorAssociation)
		case n.AssociationID == 0:
			return errors.New(ErrorAssociationID)
		case n.Association != "application" && n.Association != "issue":
			return errors.New(ErrorAssociationEntity)
		}
	}
	return nil
}

// ByID fetches a Note from the specified datastore - used for testing
func ByID(ds datastore.Datastore, id int) (Note, error) {
	return noteByID(ds, id)
}

// ByMemberID fetches all the notes linked to a Member from the specified datastore - used for testing
func ByMemberID(ds datastore.Datastore, memberID int) ([]Note, error) {
	return notesByMemberID(ds, memberID)
}

// noteByID fetches a Note record from the specified data store
func noteByID(ds datastore.Datastore, id int) (Note, error) {

	n := Note{ID: id}

	query := queries["select-note"] + " WHERE wn.id = ?"
	err := ds.MySQL.Session.QueryRow(query, id).Scan(
		&n.ID,
		&n.Type,
		&n.MemberID,
		&n.DateCreated,
		&n.DateUpdated,
		&n.DateEffective,
		&n.Content,
	)
	if err != nil {
		return n, err
	}

	n.Attachments, err = attachments(ds, n.ID)

	return n, err
}

func notesByMemberID(ds datastore.Datastore, memberID int) ([]Note, error) {

	var xn []Note

	query := queries["select-note"] + " WHERE m.id = ? ORDER BY wn.effective_on DESC"
	rows, err := ds.MySQL.Session.Query(query, memberID)
	if err != nil {
		return xn, err
	}
	defer rows.Close()

	for rows.Next() {
		n := Note{}
		rows.Scan(
			&n.ID,
			&n.Type,
			&n.MemberID,
			&n.DateCreated,
			&n.DateUpdated,
			&n.DateEffective,
			&n.Content,
		)

		var err error
		n.Attachments, err = attachments(ds, n.ID)
		if err != nil {
			return xn, nil
		}

		xn = append(xn, n)
	}

	return xn, nil
}

func attachments(ds datastore.Datastore, noteID int) ([]Attachment, error) {

	var xa []Attachment

	query := queries["select-attachments"] + " WHERE wa.wf_note_id = ?"
	rows, err := ds.MySQL.Session.Query(query, noteID)
	if err != nil {
		return xa, err
	}
	defer rows.Close()

	for rows.Next() {
		a := Attachment{}
		err := rows.Scan(&a.ID, &a.Name, &a.URL)
		if err != nil {
			return xa, err
		}
		xa = append(xa, a)
	}

	return xa, nil
}

// GetNotes fetches notes relating, optionally those that relate to
// a particular entity 'e'. An 'entity' is a value in the db that
// describes the table (entity) to which the note is linked. For example,
// a note relating to a membership title would have the value mp_title
//func (m *Member) GetNotes(entityName string, entityID string) []note.Note {
//
//	query := `SELECT
//		wn.effective_on,
//		wn.note,
//		wna.association,
//		wna.association_entity_id
//		FROM wf_note wn
//		LEFT JOIN wf_note_association wna ON wn.id = wna.wf_note_id
//		WHERE wna.member_id = ?
//		%s %s
//		ORDER BY wn.effective_on DESC`
//
//	// filter by entity name
//	s1 := ""
//	if len(entityName) > 0 {
//		s1 = " AND " + entityName + " clause here"
//	}
//
//	// Further filter by a specific entity id
//	s2 := ""
//	if len(entityID) > 0 {
//		s2 = " AND " + entityID + " clause here"
//	}
//
//	query = fmt.Sprintf(query, s1, s2)
//	fmt.Println(query)
//
//	// Get the notes relating to this title
//	n1 := note.Note{
//		ID:            123,
//		DateCreated:   "2016-01-01",
//		DateUpdated:   "2016-02-02",
//		DateEffective: "2016-03-03",
//		Content:       "This is the actual note...",
//	}
//
//	n2 := note.Note{
//		ID:            123,
//		DateCreated:   "2016-04-01",
//		DateUpdated:   "2016-05-02",
//		DateEffective: "2016-06-03",
//		Content:       "This is the second note...",
//	}
//
//	return []note.Note{n2, n1}
//}
