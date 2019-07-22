package position

import (
	"fmt"
	"github.com/cardiacsociety/web-services/internal/member"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel position report File
func ExcelReport(ds datastore.Datastore, positions []Position) (*excelize.File, error) {

	f := excel.New([]string{
		"ID",
		"Member",
		"Email",
		"Position",
		"Organisation",
		"Start",
		"End",
		"Address 1",
		"Address 2",
		"Address 3",
		"Locality",
		"State",
		"Postcode",
		"Country",
		"Email",
		"Comment",
	})

	// data rows
	for _, p := range positions {

		// If dates are bung set to an empty string
		var startDate, endDate interface{}
		if p.StartDate.Year() > 1971 { // epoch + 1
			startDate = p.StartDate
		} else {
			startDate = ""
		}
		if p.EndDate.Year() > 1971 { // epoch + 1
			endDate = p.EndDate
		} else {
			endDate = ""
		}

		// Get member record so we can access contact location
		var address = []string{"", "", ""}
		var mail member.Location
		m, err := member.ByID(ds, p.MemberID)
		if err != nil {
			f.AddError(m.ID, "Error fetching member record: "+err.Error())
		} else {
			// ContactLocationByType returns an empty struct and an error if not found
			// so can ignore error and write an empty cell
			mail, err = m.ContactLocationByDesc("mail")
			if err != nil {
				f.AddError(m.ID, "Error fetching mail address: "+err.Error())
			}
			if len(mail.Address) > 0 {
				address[0] = mail.Address[0]
			}
			if len(mail.Address) > 1 {
				address[1] = mail.Address[1]
			}
			if len(mail.Address) > 2 {
				address[2] = mail.Address[2]
			}
		}

		data := []interface{}{
			p.MemberPositionID,
			p.Member + " [" + strconv.Itoa(p.MemberID) + "]",
			p.Email,
			p.Name + " [" + strconv.Itoa(p.ID) + "]",
			p.OrganisationName + " [" + strconv.Itoa(p.OrganisationID) + "]",
			startDate,
			endDate,
			address[0],
			address[1],
			address[2],
			mail.City,
			mail.State,
			mail.Postcode,
			mail.Country,
			m.Contact.EmailPrimary,
			p.Comment,
		}

		err = f.AddRow(data)
		if err != nil {
			msg := fmt.Sprintf("AddRow() err = %s", err)
			log.Printf(msg)
			f.AddError(p.ID, msg)
		}
	}

	// style
	f.SetColWidthByHeading("Member", 18)
	f.SetColStyleByHeading("Start", excel.DateStyle)
	f.SetColWidthByHeading("Start", 18)
	f.SetColStyleByHeading("End", excel.DateStyle)
	f.SetColWidthByHeading("End", 18)

	return f.XLSX, nil
}
