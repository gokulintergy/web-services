package main

import (
	"log"
		"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/34South/envr"
	"gopkg.in/couchbase/gocb.v1"
			)

const (
	couchUser      = "admin"
	couchPass      = "password"
	bucketName     = "csanz"
)

var ds datastore.Datastore
var cb *gocb.Bucket

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
	//syncResources()
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

