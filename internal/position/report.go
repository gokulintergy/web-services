package position

import (
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
		"Comment",
	})

	// data rows
	for _, p := range positions {

		data := []interface{}{
			p.MemberPositionID,
			p.Member + " [" + strconv.Itoa(p.MemberID) + "]",
			p.Email,
			p.Name + " [" + strconv.Itoa(p.ID) + "]",
			p.OrganisationName + " [" + strconv.Itoa(p.OrganisationID) + "]",
			p.StartDate,
			p.EndDate,
			p.Comment,
		}
		err := f.AddRow(data)
		if err != nil {
			log.Printf("AddRow() err = %s\n", err)
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
