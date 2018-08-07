package main

import (
	"log"
	"fmt"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/34South/envr"
	"gopkg.in/couchbase/gocb.v1"
	"github.com/cardiacsociety/web-services/internal/generic"
	"github.com/cardiacsociety/web-services/internal/member"
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
	Type          string                    `json:"type"`
	Gender        string                    `json:"gender,omitempty"`
	PreNom        string                    `json:"preNom,omitempty"`
	FirstName     string                    `json:"firstName,omitempty"`
	LastName      string                    `json:"lastName,omitempty"`
	PostNom       string                    `json:"postNom,omitempty"`
	Email         string                    `json:"email,omitempty"`
	Email2        string                    `json:"email2,omitempty"`
	Mobile        string                    `json:"mobile,omitempty"`
	Locations     []member.Location         `json:"locations,omitempty"`
	Title         string                    `json:"title,omitempty"`
	TitleHistory  []member.MembershipTitle  `json:"titleHistory,omitempty"`
	Status        string                    `json:"status,omitempty"`
	StatusHistory []member.MembershipStatus `json:"statusHistory,omitempty"`
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

	//bucket.Manager("", "").CreatePrimaryIndex("", true, false)

	//bucket.Upsert("u:kingarthur",
	//	User{
	//		Id: "kingarthur",
	//		Email: "kingarthur@couchbase.com",
	//		Interests: []string{"Holy Grail", "African Swallows"},
	//	}, 0)
	//
	//// Get the value back
	//var inUser User

	//var brewery Brewery
	//bucket.Get("abbey_wright_brewing_valley_inn", &brewery)
	//fmt.Printf("Brewery: %v\n", brewery)
	//
	//// Use query
	//query := gocb.NewN1qlQuery("SELECT * FROM " + bucketName)
	//rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{})
	//var row interface{}
	//for rows.Next(&row) {
	//	fmt.Printf("Row: %v", row)
	//}
}

func syncMembers() {

	// get all member ids
	xi, err := generic.GetIDs(ds, "member", "")
	if err != nil {
		log.Fatalln("mysql err", err)
	}

	for _, i := range xi {
		fmt.Println("Syncing member id", i)
		m, err := member.ByID(ds, i)
		if err != nil {
			log.Fatalln("Could not get member id", i, "-", err)
		}
		md := mapMember(*m)
		id := fmt.Sprintf("%v::%v", memberIdPrefix, m.ID)
		_, err = cb.Upsert(id, md, 0)
		if err != nil {
			log.Println("Upsert error", err)
		}
	}
}

// mapMember maps member.Member to couchbase memberDoc
func mapMember(m member.Member) MemberDoc {

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

	return MemberDoc{
		Type:          "member",
		PreNom:        m.Title,
		FirstName:     m.FirstName,
		LastName:      m.LastName,
		PostNom:       m.PostNominal,
		Email:         m.Contact.EmailPrimary,
		Email2:        m.Contact.EmailSecondary,
		Mobile:        m.Contact.Mobile,
		Locations:     locations,
		Title:         title,
		TitleHistory:  titleHistory,
		Status:        status,
		StatusHistory: statusHistory,
	}
}
