package application

import (
	"log"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel application report File
func ExcelReport(ds datastore.Datastore, applications []Application) (*excelize.File, error) {

	f := excel.New([]string{
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
	f.SetColStyleByHeading("Application date", excel.DateStyle)
	f.SetColWidthByHeading("Application date", 18)

	return f.XLSX, nil
}
