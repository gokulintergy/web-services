package issue

import (
	"errors"
	"fmt"

	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Error messages
const (
	ErrorIDNotNil              = "cannot insert an issue row because ID already has a value"
	ErrorNoTypeID              = "cannot insert an issue row because Type.ID field is not set"
	ErrorNoDescription         = "cannot insert an issue row because Description is empty"
	ErrorAssociationNoMemberID = "cannot associate issue with another entity unless a member id is specified"
	ErrorAssociation           = "association entity not specified"
	ErrorAssociationID         = "association entity ID not specified"
	ErrorAssociationEntity     = "association entity unknown"
)

// Issue represents a workflow issue
type Issue struct {
	ID          int
	Type        Type
	Resolved    bool
	Visible     bool
	Description string
	Action      string
	Notes       []note.Note

	// The following fields represent data associated with an issue. In the relational database this
	// is how Issues are linked to members and invoices. The association is optionsal and open
	// so as to allow Issue to be raised without any association (gloabl issues) or specifically
	//related to a member or invoice record. As such, any connections must be determined programatically.
	// At this stage issues can only be associated with "application" or "invoice" records

	MemberID      int    // if set, this issue will be associated with this member
	Association   string // either "application" or "invoice"
	AssociationID int    // the id of the associated application or invoice record
}

// Type represents the sub-category of the issue, ie Category -> Type
type Type struct {
	ID          int
	Category    Category
	Name        string
	Description string
}

// Category represents the top-level categorisation of issues
type Category struct {
	ID          int
	Name        string
	Description string
}

// InsertRow creates a new issue row with fields from Issue
func (i *Issue) InsertRow(ds datastore.Datastore) error {
	switch {
	case i.ID > 0:
		return errors.New(ErrorIDNotNil)
	case i.Type.ID == 0:
		return errors.New(ErrorNoTypeID)
	case i.Description == "":
		return errors.New(ErrorNoDescription)
	}
	q := fmt.Sprintf(queries["insert-issue"], i.Type.ID, i.Description, i.Action)
	_, err := ds.MySQL.Session.Exec(q)
	if err != nil {
		return err
	}

	// Associate other data if fields are set, but ensure data looks ok
	err = i.checkAssociatioData()
	if err != nil {
		return err
	}

	return nil
}

// checkAssociatioData verifies fields required to associate an issue with other data
func (i *Issue) checkAssociatioData() error {
	// no association
	if i.MemberID == 0 && i.Association == "" && i.AssociationID == 0 {
		return nil
	}
	// associate with member only
	if i.MemberID > 0 && i.Association == "" && i.AssociationID == 0 {
		return nil
	}
	// associate with entity
	if i.Association != "" || i.AssociationID > 0 {
		switch {
		case i.MemberID == 0:
			return errors.New(ErrorAssociationNoMemberID)
		case i.Association == "":
			return errors.New(ErrorAssociation)
		case i.AssociationID == 0:
			return errors.New(ErrorAssociationID)
		case i.Association != "application" && i.Association != "invoice":
			return errors.New(ErrorAssociationEntity)
		}
	}
	return nil
}
