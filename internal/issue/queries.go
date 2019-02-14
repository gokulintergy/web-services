package issue

var queries = map[string]string{
	"insert-issue": insertIssue,
}

const insertIssue = `
INSERT INTO wf_issue (
	wf_issue_type_id, 
	updated_at, 
	live_on, 
	description, 
	required_action
) VALUES (%d, NOW(), NOW(), %q, %q)`
