package main

import (
	"testing"
)

var (
	ref = Reference{
		Type:    "article",
		Key:     "MP2",
		Authors: []string{"C. M{\\o}ller", "M. S. Plesset"},
		Title:   "Note on an Approximation Treatment for Many-Electron Systems",
		Journal: "Phys. Rev.",
		Volume:  "46",
		Pages:   "618-622",
		Year:    "1934",
		Tags:    []string{"test", "tag", "now"}}
)

func TestAddType(t *testing.T) {

	t.Run("Old struct but no type", func(t *testing.T) {
		start := Reference{
			Type:    "",
			Key:     "MP2",
			Authors: []string{"C. M{\\o}ller", "M. S. Plesset"},
			Title:   "Note on an Approximation Treatment for Many-Electron Systems",
			Journal: "Phys. Rev.",
			Volume:  "46",
			Pages:   "618-622",
			Year:    "1934",
			Tags:    []string{"test", "tag", "now"}}
		goaltype := "book"
		err := start.AddType(goaltype)
		if err != nil || start.Type != goaltype {
			t.Errorf("got %s, wanted %s", start.Type, goaltype)
		}
	})

	t.Run("New struct and no type", func(t *testing.T) {
		start := new(Reference)
		goaltype := "book"
		err := start.AddType(goaltype)
		if err != nil || start.Type != goaltype {
			t.Errorf("got %s, wanted %s", start.Type, goaltype)
		}
	})

	t.Run("Already have a type", func(t *testing.T) {
		goaltype := "book"
		err := ref.AddType(goaltype)
		if err == nil {
			t.Errorf("wanted an error but didn't get one")
		}
	})
}

func TestEditType(t *testing.T) {
	start := Reference{
		Type:    "article",
		Key:     "MP2",
		Authors: []string{"C. M{\\o}ller", "M. S. Plesset"},
		Title:   "Note on an Approximation Treatment for Many-Electron Systems",
		Journal: "Phys. Rev.",
		Volume:  "46",
		Pages:   "618-622",
		Year:    "1934",
		Tags:    []string{"test", "tag", "now"}}
	goaltype := "book"
	start.EditType(goaltype)
	if start.Type != goaltype {
		t.Errorf("got %s, wanted %s", start.Type, goaltype)
	}
}
