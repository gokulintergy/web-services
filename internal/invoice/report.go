package invoice

import (
	"fmt"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/platform/excel"
)

// ExcelReport returns an excel invoice report File
func ExcelReport(ds datastore.Datastore, invoices []Invoice) (*excelize.File, error) {

	f := excel.New([]string{
		"Invoice ID",
		"Invoice date",
		"Due date",
		"Member",
		"Subscription",
		"Amount",
		"Paid",
		"Comment",
	})

	// data rows
	var total float64
	for _, i := range invoices {

		paid := "no"
		if i.Paid == true {
			paid = "yes"
		}

		data := []interface{}{
			i.ID,
			i.IssueDate,
			i.DueDate,
			i.Member + " [" + strconv.Itoa(i.MemberID) + "]",
			i.Subscription,
			i.Amount,
			paid,
			i.Comment,
		}
		err := f.AddRow(data)
		if err != nil {
			msg := fmt.Sprintf("AddRow() err = %s", err)
			log.Printf(msg)
			f.AddError(i.ID, msg)
		}

		total += i.Amount
	}

	// total row
	r := []interface{}{"", "", "", "", "Total", total, "", ""}
	err := f.AddRow(r)
	if err != nil {
		msg := fmt.Sprintf("AddRow() err = %s\n", err)
		log.Printf(msg)
		f.AddError(0, msg)
	}

	// style
	f.SetColStyleByHeading("Invoice date", excel.DateStyle)
	f.SetColWidthByHeading("Invoice date", 18)
	f.SetColStyleByHeading("Due date", excel.DateStyle)
	f.SetColWidthByHeading("Due date", 18)
	f.SetColWidthByHeading("Member", 18)
	f.SetColStyleByHeading("Amount", excel.CurrencyStyle)
	f.SetColWidthByHeading("Amount", 18)
	cell := "E" + strconv.Itoa(f.NextRow)
	f.SetCellStyle(cell, cell, excel.BoldStyle)
	cell = "F" + strconv.Itoa(f.NextRow)
	f.SetCellStyle(cell, cell, excel.BoldCurrencyStyle)

	return f.XLSX, nil
}
