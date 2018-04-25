// Package main provides ...
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gorp.v1"

	"github.com/frostyfrog/urlsaver/bookmark"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/mattn/go-sqlite3"
)

var bookmarks *bookmark.BookmarkManager

func main() {
	log.Println("Starting up 地図ウエブ「ちずウエブ」")
	var db *sql.DB
	//db, err := sql.Open("sqlite3", "urlsaver_gotest2.db")
	//db, err := gorm.Open("mysql", "urlparser:@/urlparser?charset=utf8&parseTime=True&loc=Local")
	//db, err := sql.Open("mysql", "tcp:localhost:3306*urlsaver/urlsaver/urlsaver")
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	dbType := "mysql"
	if hostname == "Atlanta" {
		db, err = sql.Open("mysql", "urlsaver@/urlsaver")
	} else if hostname == "Ogre" {
		db, err = sql.Open("sqlite3", "/tmp/urlsaver_godev.db")
		dbType = "sqlite"
	} else {
		db, err = sql.Open("mysql", "urlsaver@tcp(172.24.0.2:3306)/urlsaver")
	}
	if err != nil {
		log.Fatalf("Failed to open sqlite DB: %v", err)
	}
	defer db.Close()

	var dbmap *gorp.DbMap
	if dbType == "mysql" {
		dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	} else { // Assume sqlite
		dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	}
	dbmap.AddTableWithName(bookmark.Bookmark{}, "bookmarks").SetKeys(true, "ID")
	bookmarks = bookmark.NewBookmarkManager(dbmap)
	bookmarks.SetupDB()

	r := mux.NewRouter()
	r.HandleFunc("/os", listOS).Methods("GET")
	r.HandleFunc("/bookmark", listBookmarks).Methods("GET")
	r.HandleFunc("/bookmark", newBookmark).Methods("POST")
	r.HandleFunc("/bookmark/{id}", getBookmark).Methods("GET")
	r.HandleFunc("/bookmark/{id}", updateBookmark).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Printf("Now listening on port 【:7777】\n")
	log.Fatal(http.ListenAndServe(":7777", r))
}

func listOS(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Hostname string
	}
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		log.Println(err)
	}

	// Var population
	data.Hostname = hostname

	// Response
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		log.Println(err)
	}
}

func listBookmarks(w http.ResponseWriter, r *http.Request) {
	var res struct{ bookmarks []*bookmark.Bookmark }
	var err error
	res.bookmarks, err = bookmarks.All()
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		log.Println(err)
	}
	retval := make([]bookmark.Bookmark, len(res.bookmarks))
	for k, v := range res.bookmarks {
		retval[k] = *v
	}
	err = json.NewEncoder(w).Encode(retval)
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		log.Println(err)
	}
}
func newBookmark(w http.ResponseWriter, r *http.Request) {
	req := struct{ URL string }{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Creating new bookmark '%v'", req.URL)
	t, err := bookmark.NewBookmark(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = bookmarks.Save(t)
	if err != nil {
		http.Error(w, "Failed to save bookmark", http.StatusInternalServerError)
		log.Println(err)
	}
}

func getBookmark(w http.ResponseWriter, r *http.Request) {
	txt := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(txt, 10, 0)
	if err != nil {
		http.Error(w, "bookmark ID is not a number", http.StatusBadRequest)
		return
	}
	t, ok := bookmarks.Find(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		log.Println(err)
	}
}

func updateBookmark(w http.ResponseWriter, r *http.Request) {
	txt := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(txt, 10, 0)
	if err != nil {
		http.Error(w, "task ID is not a number", http.StatusBadRequest)
		return
	}
	var t bookmark.Bookmark
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "decode failed", http.StatusBadRequest)
		return
	}
	if t.ID != id {
		http.Error(w, "inconsistent task ID", http.StatusBadRequest)
		return
	}
	if _, ok := bookmarks.Find(t.ID); !ok {
		http.NotFound(w, r)
		return
	}
	bookmarks.Save(&t)
}
