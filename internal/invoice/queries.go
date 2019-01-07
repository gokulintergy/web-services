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
<<<<<<< HEAD
    IFNULL(i.fn_subscription_id, 0) AS SubscriptionID,
    COALESCE(s.name, '') as Subscription,
=======
    i.fn_subscription_id AS SubscriptionID,
    s.name as Subscription,
>>>>>>> dab11a705fb65953e221529b4a3e75261f76c5d0
    i.start_on AS FromDate,
    i.end_on AS ToDate,
    i.invoice_total AS Amount,
    i.paid AS Paid,
<<<<<<< HEAD
    COALESCE(i.comment,'') AS Comment
=======
    i.comment AS Comment
>>>>>>> dab11a705fb65953e221529b4a3e75261f76c5d0
FROM
    fn_m_invoice i
        LEFT JOIN
    member m ON i.member_id = m.id
    LEFT JOIN fn_subscription s ON i.fn_subscription_id = s.id
WHERE 1
`

const selectActiveInvoices = selectInvoices + ` AND i.active = 1 `

const selectInvoiceByID = selectActiveInvoices + ` AND i.id = %v `
