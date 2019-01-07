package invoice

var queries = map[string]string{
	"select-invoices":      selectActiveInvoices,
	"select-invoice-by-id": selectInvoiceByID,
}

const selectInvoices = `
SELECT 
    i.id AS InvoiceID,
    i.created_at AS Created,
    i.updated_at AS Updated,
    i.member_id AS MemberID,
	COALESCE(CONCAT(m.first_name, ' ', m.last_name), '') AS Member,
	i.invoiced_on as IssueDate,
	i.last_sent_at as LastSentDate,
    i.due_on AS DueDate,
    i.fn_subscription_id AS SubscriptionID,
    s.name as Subscription,
    i.start_on AS FromDate,
    i.end_on AS ToDate,
    i.invoice_total AS Amount,
    i.paid AS Paid,
    i.comment AS Comment
FROM
    fn_m_invoice i
        LEFT JOIN
    member m ON i.member_id = m.id
    LEFT JOIN fn_subscription s ON i.fn_subscription_id = s.id
WHERE 1
`

const selectActiveInvoices = selectInvoices + ` AND i.active = 1 `

const selectInvoiceByID = selectActiveInvoices + ` AND i.id = %v `
