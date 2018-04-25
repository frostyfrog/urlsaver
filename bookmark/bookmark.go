// Package urlsaver provides ...
package bookmark

import (
	"errors"

	"github.com/badoux/goscraper"
)

type Bookmark struct {
	ID          int64  `db:"bookmark_id"`
	URL         string `db:"url"`         // URL for the bookmark
	Title       string `db:"title"`       // Title of the bookmark
	Description string `db:"description"` // Description of the bookmark
	Folder      string `db:"folder_path"` // Path to the bookmark, usually used for imported bookmarks
}

func NewBookmark(url string) (*Bookmark, error) {
	if url == "" {
		return nil, errors.New("empty URL")
	}
	scrape, err := goscraper.Scrape(url, 5)
	if err != nil {
		return nil, err
	}
	return &Bookmark{0, url, scrape.Preview.Title, scrape.Preview.Description, ""}, nil
}
