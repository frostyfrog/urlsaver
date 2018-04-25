// Package bookmark_test provides ...
package bookmark

import "testing"

func TestNewTag(t *testing.T) {
	tag := newTag("test-tag")
	if tag == nil {
		t.Fatalf("expected tag, got nil")
	}
	if tag.Name != "test-tag" {
		t.Fatalf("expected test-tag, got %v", tag.Name)
	}
}
