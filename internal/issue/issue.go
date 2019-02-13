package issue

import (
	"github.com/cardiacsociety/web-services/internal/note"
)

// Issue represents a workflow issue
type Issue struct {
	ID          int
	Type        IssueType
	Resolved    bool
	Visible     bool
	Description string
	Action      string

	Notes []note.Note
}

// IssueType represents the sub-category of the issue, ie Category -> Type
type IssueType struct {
	ID          int
	Category    IssueCategory
	Name        string
	Description string
}

// IssueCategory represents the top-level categorisation of issues
type IssueCategory struct {
	ID          int
	Name        string
	Description string
}
