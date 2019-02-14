package issue

import (
	"errors"

	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

const (
	ErrorIDNotNil = "cannot insert an issue row because ID already has a value"
	ErrorNoTypeID      = "cannot insert an issue row because Type.ID field is not set"
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

	if i.ID > 0 {
		return errors.New(ErrorIDNotNil)
	}
	if i.Type.ID == 0 {
		return errors.New(ErrorNoTypeID)
	}


	//q := fmt.Sprintf(queries["insert-issue"], i.Type.ID, i.Description, i.Action)
	return nil
}
