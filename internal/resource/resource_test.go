package resource_test

import (
	"testing"
	"time"

	"github.com/cardiacsociety/web-services/testdata"
	"github.com/cardiacsociety/web-services/internal/resource"
	"gopkg.in/mgo.v2/bson"
	"github.com/matryer/is"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

var data = testdata.NewDataStore()

func TestResources(t *testing.T) {
	is := is.New(t)
	err := data.SetupMySQL()
	is.NoErr(err) // error setting up test mysql database
	defer data.TearDownMySQL()

	err = data.SetupMongoDB()
	is.NoErr(err) // error setting up test mongo database
	defer data.TearDownMongoDB()

	t.Run("resources", func(t *testing.T) {
		t.Run("testPingDatabases", testPingDatabases)
		t.Run("testByID", testByID)
		t.Run("testDocResourcesAll", testDocResourcesAll)
		t.Run("testDocResourcesLimit", testDocResourcesLimit)
		t.Run("testDocResourcesOne", testDocResourcesOne)
		t.Run("testQueryResourcesCollection", testQueryResourcesCollection)
		t.Run("testFetchResources", testFetchResources)
		t.Run("testSyncResource", testSyncResource)
		t.Run("testSaveNewResource", testSaveNewResource)
		t.Run("testSaveExistingResource", testSaveExistingResource)
	})
}

func testPingDatabases(t *testing.T) {
	is := is.New(t)
	err := data.Store.MySQL.Session.Ping()
	is.NoErr(err) // Could not ping test mysql database
	err = data.Store.MongoDB.Session.Ping()
	is.NoErr(err) // Could not ping test mongo database
}

// tests fetching resources from mysql
func testByID(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		id  int
		doi string
	}{
		{id: 6576, doi: "https://doi.org/10.1053/j.gastro.2017.08.022"},
		{id: 6578, doi: "https://doi.org/10.1016/j.jaci.2017.03.020"},
	}

	for _, c := range cases {
		r, err := resource.ByID(data.Store, c.id)
		is.NoErr(err)                  // error fetching resource by id from mysql
		is.Equal(c.doi, r.ResourceURL) // incorrect resource url
	}
}

func testDocResourcesAll(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		query       bson.M
		projection  bson.M
		resultCount int
	}{
		{query: bson.M{}, projection: bson.M{}, resultCount: 5},
		{query: bson.M{"id": 24967}, projection: bson.M{}, resultCount: 1},
		{query: bson.M{"id": 25000}, projection: bson.M{}, resultCount: 0},
		{query: bson.M{"keywords": "PD2018"}, projection: bson.M{}, resultCount: 1},
	}

	for _, c := range cases {
		xr, err := resource.DocResourcesAll(data.Store, c.query, c.projection)
		is.NoErr(err)                    // error fetching docs
		is.Equal(c.resultCount, len(xr)) // incorrect result count
	}
}

func testDocResourcesLimit(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		query       bson.M
		projection  bson.M
		limit       int
		expectCount int
	}{
		{query: bson.M{}, projection: bson.M{}, limit: 2, expectCount: 2},
	}

	for _, c := range cases {
		xr, err := resource.DocResourcesLimit(data.Store, c.query, c.projection, c.limit)
		is.NoErr(err)                    // error fetching docs
		is.Equal(c.expectCount, len(xr)) // incorrect result count
	}
}

func testDocResourcesOne(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		id  int
		doi string
	}{
		{id: 2000, doi: "https://webcast.gigtv.com.au/Mediasite/Play/bb4663e0c3b64cc58f200064bb6c03db1d"},
		{id: 10012, doi: "https://doi.org/10.1016/j.resuscitation.2017.08.218"},
	}

	for _, c := range cases {
		r, err := resource.DocResourcesOne(data.Store, bson.M{"id": c.id})
		is.NoErr(err)                  // error fetching doc
		is.Equal(c.doi, r.ResourceURL) // incorrect resource url
	}
}

func testQueryResourcesCollection(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		query datastore.MongoQuery
		doi   string
	}{
		{
			query: datastore.MongoQuery{Find: bson.M{"id": 2000}},
			doi:   "https://webcast.gigtv.com.au/Mediasite/Play/bb4663e0c3b64cc58f200064bb6c03db1d",
		},
		{
			query: datastore.MongoQuery{Find: bson.M{"id": 10012}},
			doi:   "https://doi.org/10.1016/j.resuscitation.2017.08.218",
		},
	}

	for _, c := range cases {
		xr, err := resource.QueryResourcesCollection(data.Store, c.query)
		is.NoErr(err)        // query error
		is.Equal(len(xr), 1) // expected only one result
		r := xr[0].(bson.M)  // it returns []interface{} so need to assert
		doi, ok := r["resourceUrl"]
		if !ok {
			t.Fatal("No ResourceUrl field")
		}
		is.Equal(c.doi, doi) // incorrect resource url
	}
}
func testFetchResources(t *testing.T) {
	is := is.New(t)
	cases := []struct {
		query       map[string]interface{}
		expectCount int
	}{
		{
			query:       map[string]interface{}{"id": 2000},
			expectCount: 1,
		},
		{
			query:       map[string]interface{}{},
			expectCount: 5,
		},
	}

	for _, c := range cases {
		xr, err := resource.FetchResources(data.Store, c.query, 0) // limit = 0
		is.NoErr(err)                                              // query error
		is.Equal(len(xr), c.expectCount)                           // unexpected results count
	}
}

func testSyncResource(t *testing.T) {
	is := is.New(t)
	r, err := resource.ByID(data.Store, 6576) // this id is present in mysql but NOT in mongo test set
	is.NoErr(err)                             // error fetching resource from mysql
	resource.SyncResource(data.Store, r)      // run in a go routine so no error to check...
	time.Sleep(1 * time.Second)               // bad... but needs to time to sync to mono
	rd, err := resource.DocResourcesOne(data.Store, bson.M{"id": 6576})
	is.NoErr(err)         // error fetching sync'd resource from mongo
	is.Equal(r.ID, rd.ID) // sync'd resource has a different id?
}

func testSaveNewResource(t *testing.T) {
	is := is.New(t)

	r := resource.Resource{
		Name:        "test resource",
		ResourceURL: "http://csanz.io/abcd1234",
	}
	newId, err := r.Save(data.Store)
	is.NoErr(err) // error saving resource

	r2, err := resource.ByID(data.Store, newId)
	is.NoErr(err)                           // error fetching new resource
	is.Equal(r.Name, r2.Name)               // New resource name does not match
	is.Equal(r.ResourceURL, r2.ResourceURL) // New resource url does not match
}

func testSaveExistingResource(t *testing.T) {
	is := is.New(t)

	// this one already exists in test mysql db (id 6578) so should error as nothing to change
	r := resource.Resource{
		ResourceURL: "https://doi.org/10.1016/j.jaci.2017.03.020",
	}
	id, err := r.Save(data.Store)
	is.NoErr(err)        // save should not return an error if nothing is updated
	is.Equal(id, 6578)   // expect existing resource id to be returned
	is.Equal(r.ID, 6578) // r.ID should be set to id

	// Same again, only this time change the title
	r = resource.Resource{
		Name:        "New name",
		ResourceURL: "https://doi.org/10.1016/j.jaci.2017.03.020",
	}
	id, err = r.Save(data.Store)
	is.NoErr(err)                            // save should not return an error when record was modified
	is.Equal(id, 6578)                       // expect existing resource id to be returned
	r2, err := resource.ByID(data.Store, id) // re-fetch update resource
	is.NoErr(err)                            // error fetching the updated record
	is.Equal(r.Name, r2.Name)                // name was not updated

	//r2, err := resource.ByID(data.Store, newId)
	//is.NoErr(err)                           // error fetching new resource
	//is.Equal(r.Name, r2.Name)               // New resource name does not match
	//is.Equal(r.ResourceURL, r2.ResourceURL) // New resource url does not match
}
