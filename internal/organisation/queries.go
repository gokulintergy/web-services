package organisation

const querySelectActiveOrganisation = `SELECT 
	id, 
	short_name, 
	name, phone, 
	fax, 
	email, 
	web 
FROM 
	organisation 
WHERE 
	active = 1`

const querySelectParentOrganisations = querySelectActiveOrganisation + ` AND parent_organisation_id IS NULL`

const querySelectOrganisationByID = querySelectActiveOrganisation + ` AND id = ?`

const querySelectOrganisationByParentID = querySelectActiveOrganisation + ` AND parent_organisation_id = ?`

const querySelectOrganisationByTypeID = querySelectActiveOrganisation + ` AND organisation_type_id = ?`
