// Package urlsaver provides ...
package bookmark

import (
	"testing"
)

func newBookmarkOrFatal(t *testing.T, url string) *Bookmark {
	bookmark, err := NewBookmark(url)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	return bookmark
}

func TestNewBookmark(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com/")
	if bookmark.URL != "http://example.com/" {
		t.Errorf("expected http://example.com/, got %v", bookmark.URL)
	}
}

func TestNewBookmarkWithEmptyURL(t *testing.T) {
	_, err := NewBookmark("")
	if err == nil {
		t.Errorf("expected 'empty URL' error, got %#v", err)
	}
}

func TestAddBookmarkToCategory(t *testing.T) {
	bookmark := newBookmarkOrFatal(t, "http://example.com")
	tag := newTag("test-tag")
	tag.Link(bookmark)
	if bookmark.GetTags()[0] != tag {
		t.Fatalf("expected test-tag, got %v", bookmarks.GetTags()[0])
	}
}
