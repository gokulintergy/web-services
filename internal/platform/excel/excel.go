// Package excel create excel files
package excel

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"

	"github.com/cardiacsociety/web-services/internal/application"
	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/payment"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// MemberReport returns an excel member report File
func MemberReport(members []member.Member) (*excelize.File, error) {

	var rowNum int
	file := excelize.NewFile()

	// heading row
	rowNum++
	headingStyle, _ := file.NewStyle(`{"font": {"bold": true}}`)
	file.SetCellStyle("Sheet1", "A1", "ZZ1", headingStyle)

	xt := []string{
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
	xc := columnKeys(len(xt))
	for i := range xt {
		cell := xc[i] + strconv.Itoa(rowNum) // "A1", "A2" etc
		value := xt[i]
		file.SetCellValue("Sheet1", cell, value)
	}

	// data rows
	for _, m := range members {
		rowNum++
		row := make(map[string]interface{}, len(xc))
		keys := rowKeys(xc, rowNum)
		row[keys[0]] = strconv.Itoa(m.ID)
		row[keys[1]] = m.Title
		row[keys[2]] = m.FirstName
		row[keys[3]] = strings.Join(m.MiddleNames, " ")
		row[keys[4]] = m.LastName
		row[keys[5]] = m.PostNominal
		row[keys[6]] = m.Gender
		row[keys[7]] = m.DateOfBirth
		row[keys[8]] = m.Contact.EmailPrimary
		row[keys[9]] = m.Contact.EmailSecondary
		row[keys[10]] = m.Contact.Mobile
		row[keys[11]] = m.DateOfEntry

		if len(m.Memberships) > 0 {
			row[keys[12]] = m.Memberships[0].Title
			row[keys[13]] = m.Memberships[0].Status
		}

		row[keys[14]] = m.Country

		if len(m.Tags) > 0 {
			row[keys[15]] = strings.Join(m.Tags, ", ")
		}

		row[keys[16]] = m.JournalNumber
		row[keys[17]] = m.BpayNumber

		// ContactLocationByType returns an empty struct and an error if not found
		// so can ignore error and write an empty cell
		mailLocation, _ := m.ContactLocationByDesc("mail")
		row[keys[18]] = strings.Join(mailLocation.Address, " ")
		row[keys[19]] = mailLocation.City
		row[keys[20]] = mailLocation.State
		row[keys[21]] = mailLocation.Postcode
		row[keys[22]] = mailLocation.Country

		directoryLocation, _ := m.ContactLocationByDesc("directory")
		row[keys[23]] = strings.Join(directoryLocation.Address, " ")
		row[keys[24]] = directoryLocation.City
		row[keys[25]] = directoryLocation.State
		row[keys[26]] = directoryLocation.Postcode
		row[keys[27]] = directoryLocation.Country
		row[keys[28]] = directoryLocation.Phone
		row[keys[29]] = directoryLocation.Fax
		row[keys[30]] = directoryLocation.Email
		row[keys[31]] = directoryLocation.URL

		p1, _ := m.PositionByName("First Council Affiliation")
		row[keys[32]] = p1.OrgName
		p2, _ := m.PositionByName("Second Council Affiliation")
		row[keys[33]] = p2.OrgName
		p3, _ := m.PositionByName("Third Council Affiliation")
		row[keys[34]] = p3.OrgName

		// There can be many specialities, but generally up to 3 for the report
		// they *should* be returned in order of preference
		if len(m.Specialities) > 0 {
			row[keys[35]] = m.Specialities[0].Name
		}
		if len(m.Specialities) > 1 {
			row[keys[36]] = m.Specialities[1].Name
		}
		if len(m.Specialities) > 2 {
			row[keys[37]] = m.Specialities[2].Name
		}

		for i, t := range row {
			file.SetCellValue("Sheet1", i, t)
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
func PaymentReport(ds datastore.Datastore, payments []payment.Payment) (*excelize.File, error) {

	var rowNum int
	file := excelize.NewFile()
	file.SetColWidth("Sheet1", "B", "B", 16)
	file.SetColWidth("Sheet1", "C", "C", 20)
	file.SetColWidth("Sheet1", "D", "D", 16)

	// heading row
	rowNum++
	headingStyle, _ := file.NewStyle(`{"font": {"bold": true}}`)
	file.SetCellStyle("Sheet1", "A1", "Z1", headingStyle)

	xt := []string{
		"Payment ID",
		"Payment Date",
		"Member",
		"Payment Type",
		"Amount",
		"Invoice",
		"Comment",
	}
	xc := columnKeys(len(xt))
	for i := range xt {
		cell := xc[i] + strconv.Itoa(rowNum) // "A1", "A2" etc
		value := xt[i]
		file.SetCellValue("Sheet1", cell, value)
	}

	// data rows
	var total float64
	for _, p := range payments {
		rowNum++
		keys := rowKeys(xc, rowNum)
		row := make(map[string]interface{}, len(xc))
		row[keys[0]] = p.ID

		dateStyle, err := file.NewStyle(`{"custom_number_format": "dd mmm yyyy"}`)
		if err != nil {
			log.Printf("NewStyle() err = %s\n", err)
		}
		file.SetCellStyle("Sheet1", keys[1], keys[1], dateStyle)
		row[keys[1]] = p.Date

		row[keys[2]] = p.Member + " [" + strconv.Itoa(p.MemberID) + "]"
		row[keys[3]] = p.Type

		currencyStyle, err := file.NewStyle(`{"number_format": 8}`)
		if err != nil {
			log.Printf("NewStyle() err = %s\n", err)
		}
		file.SetCellStyle("Sheet1", keys[4], keys[4], currencyStyle)
		row[keys[4]] = p.Amount
		total += p.Amount

		var invoiceAllocations []string
		for _, i := range p.Allocations {
			invoiceAllocations = append(invoiceAllocations, strconv.Itoa(i.InvoiceID))
		}
		row[keys[5]] = strings.Join(invoiceAllocations, ", ")

		row[keys[6]] = p.Comment

		for i, t := range row {
			file.SetCellValue("Sheet1", i, t)
		}
	}

	// total
	rowNum++
	style, err := file.NewStyle(`{"font":{"bold":true}}`)
	if err != nil {
		log.Printf("NewStyle() err = %s\n", err)
	}
	cellLabel := "D" + strconv.Itoa(rowNum)
	file.SetCellStyle("Sheet1", cellLabel, cellLabel, style)
	file.SetCellValue("Sheet1", cellLabel, "Total")
	cellValue := "E" + strconv.Itoa(rowNum)
	file.SetCellStyle("Sheet1", cellValue, cellValue, style)
	file.SetCellValue("Sheet1", cellValue, total)

	return file, nil
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
