// Package membership provides functionality to handle membership and
// membership-related data such as applications for membership.
package membership

import (
	"github.com/cardiacsociety/web-services/internal/member"
)

// MemberRoleId is a flag in the relational database that represents a level 
// of access for a member - full or partial. 
const MemberRoleId = 1

// Application represents an application for membership.
type Application struct {

	// Member row contains the basic member data
	member.Row

	// member record - generates id for subsequent records
    // "gender": "Male",
    // "title": "Prof",
    // "firstName": "Mike",
    // "middleNames": "",
    // "lastName": "Donnici",
    // "dateOfBirth": "1970-11-03",
    // "primaryEmail": "michael@mesa.net.au",
    // "secondaryEmail": "michael@mesa.net.au",
    // "mobile": "+61402400191",
    // "qualificationInfo": "", // qualification_other?
    // "consentDirectory": true,
    // "consentContact": true,

    // IF TRUE - mp_m_tag with mp_tag_id = 4 (Advanced Trainee)
    "trainee": true,

	// ms_m_application record
	"type": "Associate",
    "nominator": {
        "id": null,
        "name": ""
    },
    "seconder": {
        "id": null,
        "name": ""
    },
    "nominatorInfo": "ghggh", // comment

    // mp_m_qualification record
    "qualifications": [{
        "qualificationId": 2,
        "name": "Bachelor of Medicine, Bachelor of Surgery",
        "abbreviation": "MBBS",
        "year": 2001,
        "organisationId": 239,
        "organisationName": "University of Sydney, RPA Campus"
    }],

    // mp_m_specialities
    "interests": ["Electrophysiology and Pacing", 1, {
        "id": 1,
        "name": "Cardiac Care Nurse (Medical)"
    }, {
        "id": 2,
        "name": "Cardiac Cath Lab Nurse"
    }, {
        "id": 3,
        "name": "Cardiac Technologist"
    }, {
        "id": 41,
        "name": "Statistician"
    }, {
        "id": 39,
        "name": "Rehab Nurse"
    }, {
        "id": 36,
        "name": "Physiotherapist"
    }],

    // mp_m_position with mp_position_id = 1 (First Council), 2 (Second Council) or 3 (Third Council)
    // Note: may have to rework this as currently not limited to 3, and no order is enforced
    "councils": ["General Cardiology", {
        "id": 1,
        "code": "CL_AH",
        "name": "Allied Health Council"
    }, {
        "id": 2,
        "code": "CL_IMAGING",
        "name": "Cardiac Imaging Council"
    }, {
        "id": 3,
        "code": "CL_SURGERY",
        "name": "Cardiovascular Surgery Council"
    }, {
        "id": 322,
        "code": "AHST Council",
        "name": "Allied Health Science and Technology Council"
    }, {
        "id": 135,
        "code": "CL_CC",
        "name": "Coronary Care Council"
    }, {
        "id": 134,
        "code": "CL_BASICSC",
        "name": "Basic Science Council"
    }],

    // Into note - other possibilities for a note might be:
    // acknowledge privacy policy, agree to abide by constitution
    "consentRequestInfo": true,

    // Workflow item?
    "ishr": true
}
