package bookmark

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

func TestSaveBookmark(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")

	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	err := m.Save(bookmark)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
}

func TestSaveBookmarkAndRetrieve(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")

	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	err := m.Save(bookmark)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	all, err := m.All()
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected one bookmark, got %v", len(all))
		if all[0].URL != bookmark.URL {
			t.Errorf("expected title %q, got %q", bookmark.URL, all[0].URL)
		}
	}
}

func TestSaveAndRetrieveTwoBookmarks(t *testing.T) {
	learnGo := newBookmarkOrFatal(t, "http://example.com/")
	learnTDD := newBookmarkOrFatal(t, "http://subdomain.example.com/")

	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	m.Save(learnGo)
	m.Save(learnTDD)

	all, err := m.All()
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 bookmarks, got %v", len(all))
	}
	if *all[0] != *learnGo {
		t.Errorf("missing bookmark: %v", learnGo)
	}
	if *all[1] != *learnTDD {
		t.Errorf("missing bookmark: %v", learnTDD)
	}
}

func TestSaveModifyAndRetrieve(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")
	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	m.Save(bookmark)

	bookmark.Title = "Example site"

	all, err := m.All()
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if all[0].Title == bookmark.Title {
		t.Errorf("saved bookmark wasn't done")
	}
}

func TestSaveTwiceAndRetrieve(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")
	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	m.Save(bookmark)
	m.Save(bookmark)

	all, err := m.All()
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 bookmark, got %v", len(all))
	}
	if *all[0] != *bookmark {
		t.Errorf("expected bookmark %v, got %v", bookmark, all[0])
	}
}
func TestSaveModifySaveAndRetrieve(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")
	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	m.Save(bookmark)
	bookmark.Title = "Example site"
	m.Save(bookmark)

	all, err := m.All()
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 bookmark, got %v", len(all))
	}
	if *all[0] != *bookmark {
		t.Errorf("expected bookmark %v, got %v", bookmark, all[0])
	}
}

func TestSaveAndFind(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")
	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)
	m.SetupDB()
	m.Save(bookmark)

	nb, ok := m.Find(bookmark.ID)
	if !ok {
		t.Errorf("Didn't find taswk")
	}
	if *bookmark != *nb {
		t.Errorf("Expected %v, got %v", bookmark, nb)
	}
}

func TestDBCreationFromStruct(t *testing.T) {
	db, err := sql.Open("sqlite3", "/tmp/urlsaver_gotest.db")
	checkError(t, err, "Failed to open sqlite DB")
	defer db.Close()

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Bookmark{}, "bookmarks").SetKeys(true, "ID")

	err = dbmap.CreateTables()
	checkError(t, err, "Create tables failed")
	err = dbmap.DropTables()
	checkError(t, err, "Drop tables failed")
}

func checkError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf(msg+": %v", err)
	}
}

func createDBMap(t *testing.T, path string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", "/tmp/urlsaver_gotest.db")
	if err != nil {
		db.Close()
		t.Fatalf("Failed to open sqlite DB: %v", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Bookmark{}, "bookmarks").SetKeys(true, "ID")
	return dbmap
}

func dbCleanup(db *gorp.DbMap) {
	db.DropTables()
	db.Db.Close()
}

func TestBookmarkManagerDBSetup(t *testing.T) {
	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)

	// Setup the DB twice, first successfully, second as well
	err := m.SetupDB()
	checkError(t, err, "Failed to setup DB")
	err = m.SetupDB()
	checkError(t, err, "Failed to setup DB (second time)")

	// Check to see if the DB is set
	if m.Db != db {
		t.Errorf("Expected db to be set")
	}

	// Create a bookmark to add to the DB to test the structure
	bm := newBookmarkOrFatal(t, "http://example.com/")
	bm.Title = "Example Site"
	bm.ID = 42

	// Make sure that the schema exists and is correct
	err = db.Insert(bm)
	checkError(t, err, "Failed to insert dummy item")

	emptyBookmark := Bookmark{}
	err = db.SelectOne(&emptyBookmark, "select * from bookmarks")
	checkError(t, err, "Failed to select dummy item")
	if emptyBookmark != *bm {
		t.Errorf("Expected %v, got %v in select", *bm, emptyBookmark)
	}
}

func TestBookmarkManagerDBCommit(t *testing.T) {
	db := createDBMap(t, "/tmp/urlsaver_gotest.db")
	defer dbCleanup(db)

	m := NewBookmarkManager(db)

	// Setup the DB twice, first successfully, second as well
	err := m.SetupDB()
	checkError(t, err, "Failed to setup DB")

	// Create a bookmark to add to the DB to test the structure
	bm := newBookmarkOrFatal(t, "http://example.com/")
	bm.Title = "Example Site"
	m.Save(bm)

	emptyBookmark := Bookmark{}
	err = db.SelectOne(&emptyBookmark, "select * from bookmarks")
	checkError(t, err, "Failed to select dummy item")
	if emptyBookmark.ID != 1 {
		t.Errorf("expected ID 1, got %v", emptyBookmark.ID)
	}
	bm.ID = 1
	if emptyBookmark != *bm {
		t.Errorf("expected %v, got %v", *bm, emptyBookmark)
	}
}
