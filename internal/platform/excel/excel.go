// Package excel create excel files
package excel

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"

	"github.com/cardiacsociety/web-services/internal/application"
	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/payment"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// MemberReport returns an excel member report File
func MemberReport(members []member.Member) (*xlsx.File, error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("file.AddSheet() err = %s", err)
	}

	columns := []string{
		"Member ID",
		"Prefix",
		"First Name",
		"Middle Name(s)",
		"Last Name",
		"Suffix",
		"Gender",
		"Date of Birth",
		"Email (primary)",
		"Email (secondary)",
		"Mobile",
		"Date of Entry",
		"Membership Title",
		"Membership Status",
		"Membership Country",
		"Tags",
		"Journal No.",
		"BPAY No.",
		"Mail Address",
		"Mail Locality",
		"Mail State",
		"Mail Postcode",
		"Mail Country",
		"Directory Address",
		"Directory Locality",
		"Directory State",
		"Directory Postcode",
		"Directory Country",
		"Directory Phone",
		"Directory Fax",
		"Directory Email",
		"Directory Web",
		"First Council",
		"Second Council",
		"Third Council",
		"First Speciality",
		"Second Speciality",
		"Third Speciality",
	}

	// Column headers
	row = sheet.AddRow()
	for _, c := range columns {
		cell = row.AddCell()
		cell.Value = c
	}

	// member rows
	for _, m := range members {
		row = sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(m.ID)
		row.AddCell().Value = m.Title
		row.AddCell().Value = m.FirstName
		row.AddCell().Value = strings.Join(m.MiddleNames, " ")
		row.AddCell().Value = m.LastName
		row.AddCell().Value = m.PostNominal
		row.AddCell().Value = m.Gender
		row.AddCell().Value = m.DateOfBirth
		row.AddCell().Value = m.Contact.EmailPrimary
		row.AddCell().Value = m.Contact.EmailSecondary
		row.AddCell().Value = m.Contact.Mobile
		row.AddCell().Value = m.DateOfEntry

		if len(m.Memberships) > 0 {
			row.AddCell().Value = m.Memberships[0].Title
			row.AddCell().Value = m.Memberships[0].Status
		}

		row.AddCell().Value = m.Country

		if len(m.Tags) > 0 {
			row.AddCell().Value = strings.Join(m.Tags, ", ")
		}

		row.AddCell().Value = m.JournalNumber
		row.AddCell().Value = m.BpayNumber

		// ContactLocationByType returns an empty struct and an error if not found
		// so can ignore error and write an empty cell
		mailLocation, _ := m.ContactLocationByDesc("mail")
		row.AddCell().Value = strings.Join(mailLocation.Address, " ")
		row.AddCell().Value = mailLocation.City
		row.AddCell().Value = mailLocation.State
		row.AddCell().Value = mailLocation.Postcode
		row.AddCell().Value = mailLocation.Country

		directoryLocation, _ := m.ContactLocationByDesc("directory")
		row.AddCell().Value = strings.Join(directoryLocation.Address, " ")
		row.AddCell().Value = directoryLocation.City
		row.AddCell().Value = directoryLocation.State
		row.AddCell().Value = directoryLocation.Postcode
		row.AddCell().Value = directoryLocation.Country
		row.AddCell().Value = directoryLocation.Phone
		row.AddCell().Value = directoryLocation.Fax
		row.AddCell().Value = directoryLocation.Email
		row.AddCell().Value = directoryLocation.URL

		p1, _ := m.PositionByName("First Council Affiliation")
		row.AddCell().Value = p1.OrgName
		p2, _ := m.PositionByName("Second Council Affiliation")
		row.AddCell().Value = p2.OrgName
		p3, _ := m.PositionByName("Third Council Affiliation")
		row.AddCell().Value = p3.OrgName

		// There can be many specialities, but generally up to 3 for the report
		// they *should* be returned in order of preference
		if len(m.Specialities) > 0 {
			row.AddCell().Value = m.Specialities[0].Name
		}
		if len(m.Specialities) > 1 {
			row.AddCell().Value = m.Specialities[1].Name
		}
		if len(m.Specialities) > 2 {
			row.AddCell().Value = m.Specialities[2].Name
		}
	}

	return file, nil
}

// ApplicationReport returns an excel application report File
func ApplicationReport(ds datastore.Datastore, applications []application.Application) (*xlsx.File, error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("file.AddSheet() err = %s", err)
	}

	columns := []string{
		"Application ID",
		"Application date",
		"Member ID",
		"Member name",
		"Nominator ID",
		"Nominator name",
		"Seconder ID",
		"Seconder name",
		"Applied for",
		"Tags",
		"Region",
		"Result",
		"Comment",
	}

	// Column headers
	row = sheet.AddRow()
	for _, c := range columns {
		cell = row.AddCell()
		cell.Value = c
	}

	for _, a := range applications {
		row = sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(a.ID)
		row.AddCell().Value = a.Date.Format("2006-01-02")
		row.AddCell().Value = strconv.Itoa(a.MemberID)
		row.AddCell().Value = a.Member

		var nomID string
		if a.NominatorID.Int64 > 0 {
			nomID = strconv.FormatInt(a.NominatorID.Int64, 10)
		}
		row.AddCell().Value = nomID
		row.AddCell().Value = a.Nominator

		var secID string
		if a.SeconderID.Int64 > 0 {
			secID = strconv.FormatInt(a.SeconderID.Int64, 10)
		}
		row.AddCell().Value = secID
		row.AddCell().Value = a.Seconder

		row.AddCell().Value = a.For

		var tags string
		var region string
		m, err := member.ByID(ds, a.MemberID)
		if err != nil {
			log.Printf("member.ByID() err = %s", err)
			tags, region = "err", "err"
		} else {
			tags = strings.Join(m.Tags, ", ")
			region = m.Country + " " + m.Contact.Locations[0].State + " " + m.Contact.Locations[0].City
		}
		row.AddCell().Value = tags
		row.AddCell().Value = region

		var status string
		if a.Status == -1 {
			status = "pending"
		}
		if a.Status == 0 {
			status = "rejected"
		}
		if a.Status == 1 {
			status = "accepted"
		}
		row.AddCell().Value = status

		row.AddCell().Value = a.Comment
	}

	return file, nil
}

// PaymentReport returns an excel payment report File
func PaymentReport(ds datastore.Datastore, payments []payment.Payment) (*xlsx.File, error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("file.AddSheet() err = %s", err)
	}

	columns := []string{
		"Payment ID",
		"Payment Date",
		"Member ID",
		"Member",
		"Payment Type",
		"Amount",
		"Comment",
		"Field1",
		"Field2",
		"Field3",
		"Field4",
		"Invoice IDs",
	}

	// Column headers
	row = sheet.AddRow()
	for _, c := range columns {
		cell = row.AddCell()
		cell.Value = c
	}

	for _, p := range payments {
		row = sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(p.ID)
		row.AddCell().Value = p.Date.Format("2006-01-02")
		row.AddCell().Value = strconv.Itoa(p.MemberID)
		row.AddCell().Value = p.Member
		row.AddCell().Value = p.Type
		row.AddCell().Value = strconv.FormatFloat(p.Amount, 'e', -1, 64)
		row.AddCell().Value = p.Comment
		row.AddCell().Value = p.DataField1
		row.AddCell().Value = p.DataField2
		row.AddCell().Value = p.DataField3
		row.AddCell().Value = p.DataField4

		var invoiceAllocations []string
		for _, i := range p.Allocations {
			invoiceAllocations = append(invoiceAllocations, strconv.Itoa(i.InvoiceID))
		}
		row.AddCell().Value = strings.Join(invoiceAllocations, ", ")
	}

	return file, nil
}
