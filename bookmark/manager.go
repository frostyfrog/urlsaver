package bookmark

import (
	"errors"
	"log"

	"gopkg.in/gorp.v1"
)

type BookmarkManager struct {
	Db *gorp.DbMap
}

func NewBookmarkManager(db *gorp.DbMap) *BookmarkManager {
	return &BookmarkManager{Db: db}
}

func (bm *BookmarkManager) SetupDB() error {
	table := bm.Db.AddTableWithName(Bookmark{}, "bookmarks")
	url_field := table.ColMap("url")
	url_field.SetNotNull(true)
	url_field.SetUnique(true)
	err := bm.Db.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

func (tm *BookmarkManager) Save(bookmark *Bookmark) error {
	if tm.Db != nil {
		//aival, err := tm.Db.SelectInt("select seq from SQLITE_SEQUENCE where name='bookmarks'")
		//if err != nil {
		//	return err
		//}
		all, err := tm.All()
		if err != nil {
			return err
		}
		for _, t := range all {
			if t.URL == bookmark.URL {
				bookmark.ID = t.ID
				tm.Db.Update(bookmark)
				return nil
			}
		}
		tm.Db.Insert(bookmark)
		//bookmark.ID = aival + 1
		return nil
	}
	return errors.New("No DB set up")
}

func (tm *BookmarkManager) All() ([]*Bookmark, error) {
	var bookmarks []*Bookmark
	_, err := tm.Db.Select(&bookmarks, "select * from bookmarks")
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func cloneBookmark(t *Bookmark) *Bookmark {
	c := *t
	return &c
}

func (tm *BookmarkManager) Find(ID int64) (*Bookmark, bool) {
	all, err := tm.All()
	if err != nil {
		log.Printf("An error occured while trying to find ID%v: %v", ID, err)
		return nil, false
	}
	for _, t := range all {
		if (*t).ID == ID {
			return t, true
		}
	}
	return nil, false
}
