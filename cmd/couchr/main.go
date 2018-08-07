package main

import (
	"log"
	"fmt"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/34South/envr"
	"gopkg.in/couchbase/gocb.v1"
	"github.com/cardiacsociety/web-services/internal/generic"
	"github.com/cardiacsociety/web-services/internal/member"
	"time"
	"github.com/cardiacsociety/web-services/internal/cpd"
)

const (
	couchUser      = "admin"
	couchPass      = "password"
	bucketName     = "csanz"
	memberIdPrefix = "member"
)

var ds datastore.Datastore
var cb *gocb.Bucket

// MemberDoc stores a couchbase member doc
type MemberDoc struct {
	Type           string                    `json:"type"`
	Created        time.Time                 `json:"created"`
	Updated        time.Time                 `json:"updated"`
	Status         string                    `json:"status,omitempty"`
	Title          string                    `json:"title,omitempty"`
	Gender         string                    `json:"gender,omitempty"`
	PreNom         string                    `json:"preNom,omitempty"`
	FirstName      string                    `json:"firstName,omitempty"`
	MiddleNames    []string                  `json:"middleNames,omitempty"`
	LastName       string                    `json:"lastName,omitempty"`
	PostNom        string                    `json:"postNom,omitempty"`
	Email          string                    `json:"email,omitempty"`
	Email2         string                    `json:"email2,omitempty"`
	Mobile         string                    `json:"mobile,omitempty"`
	Directory      bool                      `json:"directoryConsent"`
	Consent        bool                      `json:"contactConsent"`
	Locations      []member.Location         `json:"locations,omitempty"`
	Qualifications []member.Qualification    `json:"qualifications,omitempty"`
	Accreditations []member.Accreditation    `json:"accreditations,omitempty"`
	Positions      []member.Position         `json:"positions,omitempty"`
	Specialities   []member.Speciality       `json:"specialities"`
	TitleHistory   []member.MembershipTitle  `json:"titleHistory,omitempty"`
	StatusHistory  []member.MembershipStatus `json:"statusHistory,omitempty"`
	CPD            []CPD                     `json:"cpd,omitempty"`
}

type CPD struct {
	Date       string  `json:"date"`
	Category   string  `json:"category"`
	Activity   string  `json:"activity"`
	Type       string  `json:"type"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	UnitCredit float64 `json:"creditPerUnit"`
	Credit     float64 `json:"credit"`
}

func init() {
	envr.New("couchrEnv", []string{
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_URL",
	}).Auto()
}

func main() {
	log.Println("Migrating data to CouchDB...")

	log.Println("Setting up data store from env...")

	connectDatastore()
	connectCouchDB()
	syncMembers()
}

// connect the global datastore
func connectDatastore() {
	var err error
	ds, err = datastore.FromEnv()
	if err != nil {
		log.Fatalln("Could not set datastore -", err)
	}
}

// connect the global couchbase bucket
func connectCouchDB() {

	cluster, err := gocb.Connect("couchbase://localhost")
	if err != nil {
		log.Fatalln("Could not connect to couchbase", err)
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: couchUser,
		Password: couchPass,
	})

	cb, err = cluster.OpenBucket(bucketName, "")
	if err != nil {
		log.Fatalln("Could not get bucket", err)
	}
}

func syncMembers() {

	xi, err := generic.GetIDs(ds, "member", "")
	if err != nil {
		log.Fatalln("mysql err", err)
	}

	for _, id := range xi {

		md := &MemberDoc{}

		fmt.Print("Syncing member id ", id)
		m, err := member.ByID(ds, id)
		if err != nil {
			log.Fatalln("Could not get member id ", id, "-", err)
		}
		md.mapMemberProfile(*m)

		fmt.Print("... fetching cpd\n")

		xa, err := cpd.ByMemberID(ds, id)
		if err != nil {
			log.Fatalln("Could not get CPD for member id", id, "-", err)
		}
		if len(xa) > 0 {
			md.mapCPD(xa)
		}

		id := fmt.Sprintf("%v::%v", memberIdPrefix, m.ID)
		_, err = cb.Upsert(id, md, 0)
		if err != nil {
			log.Println("Upsert error", err)
		}
	}
}

// mapMemberProfile maps profile data from member.Member to couchbase memberDoc
func (md *MemberDoc) mapMemberProfile(m member.Member) {

	var title string
	var titleHistory []member.MembershipTitle

	var status string
	var statusHistory []member.MembershipStatus

	if len(m.Memberships) > 0 {

		title = m.Memberships[0].Title
		xt := m.Memberships[0].TitleHistory
		for _, t := range xt {
			titleHistory = append(titleHistory, t)
		}

		status = m.Memberships[0].Status
		xs := m.Memberships[0].StatusHistory
		for _, s := range xs {
			statusHistory = append(statusHistory, s)
		}
	}

	var locations []member.Location
	if len(m.Contact.Locations) > 0 {
		for _, l := range m.Contact.Locations {
			locations = append(locations, l)
		}
	}

	md.Type = "member"
	md.Created = m.CreatedAt
	md.Updated = m.UpdatedAt
	md.Gender = m.Gender
	md.PreNom = m.Title
	md.FirstName = m.FirstName
	md.MiddleNames = m.MiddleNames
	md.LastName = m.LastName
	md.PostNom = m.PostNominal
	md.Email = m.Contact.EmailPrimary
	md.Email2 = m.Contact.EmailSecondary
	md.Mobile = m.Contact.Mobile
	md.Directory = m.Contact.Directory
	md.Consent = m.Contact.Consent
	md.Locations = locations
	md.Title = title
	md.TitleHistory = titleHistory
	md.Status = status
	md.StatusHistory = statusHistory
	md.Qualifications = m.Qualifications
	md.Accreditations = m.Accreditations
	md.Specialities = m.Specialities
	md.Positions = m.Positions
}

// mapCPD maps cpd.CPD values to local, simpler version
func (md *MemberDoc) mapCPD(cpd []cpd.CPD) {

	for _, c := range cpd {
		ca := CPD{}
		ca.Date = c.Date
		ca.Category = c.Category.Name
		ca.Activity = c.Activity.Name
		ca.Type = c.Type.Name
		ca.Quantity = c.CreditData.Quantity
		ca.UnitCredit = c.CreditData.UnitCredit
		ca.Unit = c.CreditData.UnitName
		ca.Credit = c.Credit

		md.CPD = append(md.CPD, ca)
	}
}
