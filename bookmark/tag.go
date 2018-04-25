// Package bookmark provides ...
package bookmark

type Tag struct {
	Name string
}

func newTag(name string) *Tag {
	return &Tag{Name: name}
}
