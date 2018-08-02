package resource_test

import (
	"testing"

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
