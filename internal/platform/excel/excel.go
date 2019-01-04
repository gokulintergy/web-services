// Package excel create excel files
package excel

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/cardiacsociety/web-services/internal/application"
	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/payment"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// MemberReport returns an excel member report File
func MemberReport(members []member.Member) (*excelize.File, error) {

	f := New([]string{
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
	})

	// data rows
	for _, m := range members {

		var dob interface{}
		dob = m.DateOfBirth // string
		d, err := time.Parse("2006-01-02", m.DateOfBirth)
		if err == nil {
			dob = d // time.Time will accept the dateStyle formatting
		}

		var doe interface{}
		doe = m.DateOfEntry // string
		de, err := time.Parse("2006-01-02", m.DateOfEntry)
		if err == nil {
			doe = de // time.Time will accept the dateStyle formatting
		}

		var title string
		var status string
		if len(m.Memberships) > 0 {
			title = m.Memberships[0].Title
			status = m.Memberships[0].Status
		}

		var tags string
		if len(m.Tags) > 0 {
			tags = strings.Join(m.Tags, ", ")
		}

		// ContactLocationByType returns an empty struct and an error if not found
		// so can ignore error and write an empty cell
		mail, _ := m.ContactLocationByDesc("mail")
		directory, _ := m.ContactLocationByDesc("directory")

		p1, _ := m.PositionByName("First Council Affiliation")
		p2, _ := m.PositionByName("Second Council Affiliation")
		p3, _ := m.PositionByName("Third Council Affiliation")

		// There can be many specialities, but generally up to 3 for the report
		// they *should* be returned in order of preference
		var s1, s2, s3 string
		if len(m.Specialities) > 0 {
			s1 = m.Specialities[0].Name
		}
		if len(m.Specialities) > 1 {
			s2 = m.Specialities[1].Name
		}
		if len(m.Specialities) > 2 {
			s3 = m.Specialities[2].Name
		}

		data := []interface{}{
			m.ID,
			m.Title,
			m.FirstName,
			strings.Join(m.MiddleNames, " "),
			m.LastName,
			m.PostNominal,
			m.Gender,
			dob,
			m.Contact.EmailPrimary,
			m.Contact.EmailSecondary,
			m.Contact.Mobile,
			doe,
			title,
			status,
			m.Country,
			tags,
			m.JournalNumber,
			m.BpayNumber,
			strings.Join(mail.Address, " "),
			mail.City,
			mail.State,
			mail.Postcode,
			mail.Country,
			strings.Join(directory.Address, " "),
			directory.City,
			directory.State,
			directory.Postcode,
			directory.Country,
			directory.Phone,
			directory.Fax,
			directory.Email,
			directory.URL,
			p1.OrgName,
			p2.OrgName,
			p3.OrgName,
			s1,
			s2,
			s3,
		}
		f.AddRow(data)
	}

	return f.XLSX, nil
}

// ApplicationReport returns an excel application report File
func ApplicationReport(ds datastore.Datastore, applications []application.Application) (*excelize.File, error) {

	f := New([]string{
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
	})

	// data rows
	for _, a := range applications {

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

		data := []interface{}{
			a.ID,
			a.Date,
			a.MemberID,
			a.Member,
			a.NominatorID.Int64,
			a.Nominator,
			a.SeconderID.Int64,
			a.Seconder,
			a.For,
			tags,
			region,
			status,
			a.Comment,
		}

		f.AddRow(data)
	}

	// customise style
	f.SetColStyleByHeading("Application date", styleDate)
	f.SetColWidthByHeading("Application date", 18)

	return f.XLSX, nil
}

// PaymentReport returns an excel payment report File
func PaymentReport(ds datastore.Datastore, payments []payment.Payment) (*excelize.File, error) {

	f := New([]string{
		"Payment ID",
		"Payment date",
		"Member",
		"Payment type",
		"Amount",
		"Invoice",
		"Comment",
	})

	// data rows
	var total float64
	for _, p := range payments {

		var ia []string
		for _, i := range p.Allocations {
			ia = append(ia, strconv.Itoa(i.InvoiceID))
		}
		invoiceAllocations := strings.Join(ia, ", ")

		data := []interface{}{
			p.ID,
			p.Date,
			p.Member + " [" + strconv.Itoa(p.MemberID) + "]",
			p.Type,
			p.Amount,
			invoiceAllocations,
			p.Comment,
		}
		f.AddRow(data)

		total += p.Amount
	}

	// total row
	r := []interface{}{"", "", "", "Total", total, "", ""}
	f.AddRow(r)

	// style
	f.SetColStyleByHeading("Payment date", styleDate)
	f.SetColWidthByHeading("Payment date", 18)
	f.SetColWidthByHeading("Member", 18)
	f.SetColStyleByHeading("Amount", styleCurrency)
	f.SetColWidthByHeading("Amount", 18)
	cell := "D" + strconv.Itoa(f.nextRow)
	f.SetCellStyle(cell, cell, styleBold)
	cell = "E" + strconv.Itoa(f.nextRow)
	f.SetCellStyle(cell, cell, styleBoldCurrency)

	return f.XLSX, nil
}

// columnkeys generates the specified number of column references - eg "A", "B" ... "Z", "AA", "AB" etc.
func columnKeys(numCols int) []string {

	result := []string{}
	xa := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	for i := 0; i < numCols; i++ {

		var colName string
		var colPrefix string

		set := int(math.Floor(float64(i) / float64(26)))
		if set > 0 {
			colPrefix = xa[set-1]
		}
		colName = colPrefix + xa[i-(set*26)]
		result = append(result, colName)
	}

	return result
}

// rowKeys returns a []string containing the cell references for a row, eg "A10", "B10", "C10" etc
func rowKeys(columnKeys []string, rowNumber int) []string {
	var refs []string
	rowNum := strconv.Itoa(rowNumber)
	for _, c := range columnKeys {
		r := c + rowNum
		refs = append(refs, r)
	}
	return refs
}
