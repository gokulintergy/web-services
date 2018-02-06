/*
	Package fileset provides information about stored file resources
*/
package fileset

import (
	"fmt"

	"database/sql"

	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// FileSet represent a row from the fs_set table and describes a "set" of related files in cloud storage.
type FileSet struct {
	// ID - fs_set id value
	ID     int    `json:"id"`
	// Entity - db table containing the records to which files in this set are 'attached'
	Entity string `json:"entity"`
	// Volume - the name of the bucket - eg 'notes'
	Volume string `json:"volume"`
	// Path - or pseudo path (S3), excluding the file name, eg '/cpd/', '/note/'
	Path   string `json:"path"`
}

// NoteAttachment returns a pointer to a FileSet with relevant values for a Note attachment
func NoteAttachment() (FileSet, error) {
	return get("wf_attachment")
}

// ActivityAttachment returns a pointer to a FileSet with relevant values for an Activity attachment
func ActivityAttachment() (FileSet, error) {
	return get("ce_m_activity_attachment")
}

// ResourceAttachment returns a pointer to a FileSet with relevant values for a Resource attachment
func ResourceAttachment() (FileSet, error) {
	return get("ol_resource_file")
}

// New returns a pointer to an initialised FileSet value. It receives the setPath, eg '/notes/' which is the base
// path / pseudo path (S3) for all files stored in the set.
func get(entity string) (FileSet, error) {

	var fs FileSet
	fs.Entity = entity

	query := "SELECT id, volume_name, set_path FROM fs_set WHERE active = 1 AND current = 1 and entity_name = ?"
	err := datastore.MySQL.Session.QueryRow(query, fs.Entity).Scan(
		&fs.ID,
		&fs.Volume,
		&fs.Path,
	)
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("No file set found for set path: %s", fs.Path)
		return fs, errors.New(msg)
	}
	if err != nil {
		return fs, errors.New("New() database error - " + err.Error())
	}

	return fs, nil
}

// CheckFields is a convenience function for validating the FileSet. It returns an error for the first field
// with a zero value.
func (fs FileSet) CheckFields() error {

	if fs.ID == 0 {
		return errors.New("FileSet.ID (int) has a zero value")
	}
	if fs.Path == "" {
		return errors.New("FileSet.Path is an empty string")
	}
	if fs.Entity == "" {
		return errors.New("FileSet.Entity is an empty string")
	}
	if fs.Volume == "" {
		return errors.New("FileSet.Volume is an empty string")
	}

	return nil
}
