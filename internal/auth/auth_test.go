package auth_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/auth"
	//"github.com/cardiacsociety/web-services/internal/auth"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestAuth(t *testing.T) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	t.Run("auth", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testAuthMemberClearPass", testAuthMemberClearPass)
		t.Run("testAuthMemberMD5Pass", testAuthMemberMD5Pass)
		t.Run("testAuthMemberFail", testAuthMemberFail)
		t.Run("testAuthAdminClearPass", testAuthAdminClearPass)
		t.Run("testAuthAdminMD5Pass", testAuthAdminMD5Pass)
		t.Run("testAuthAdminFail", testAuthAdminFail)
	})
}

func testPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func testAuthMemberClearPass(t *testing.T) {
	id, name, err := auth.AuthMember(db.Store, "michael@mesa.net.au", "password")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Michael Donnici", name)
}

func testAuthMemberMD5Pass(t *testing.T) {
	id, name, err := auth.AuthMember(db.Store, "michael@mesa.net.au", "5f4dcc3b5aa765d61d8327deb882cf99")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Michael Donnici", name)
}

func testAuthMemberFail(t *testing.T) {
	id, _, err := auth.AuthMember(db.Store, "michael@mesa.net.au", "wrongPassword")
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, sql.ErrNoRows, err)
	helper.Result(t, 0, id)
}

func testAuthAdminClearPass(t *testing.T) {
	id, name, err := auth.AdminAuth(db.Store, "demo-admin", "demo-admin")
	if err == sql.ErrNoRows {
		t.Log("Expected result to fail login")
	}
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Demo Admin", name)
}

func testAuthAdminMD5Pass(t *testing.T) {
	id, name, err := auth.AdminAuth(db.Store, "demo-admin", "41d0510a9067999b72f38ba0ce9f6195")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Demo Admin", name)
}

func testAuthAdminFail(t *testing.T) {
	id, _, err := auth.AdminAuth(db.Store, "demo-admin", "wrongPassword")
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, sql.ErrNoRows, err)
	helper.Result(t, 0, id)
}
