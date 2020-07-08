package main

import (
	"reflect"
	"testing"
)

var (
	ref = Reference{
		"article",
		"MP2",
		"C. M{\\o}ller and M. S. Plesset",
		"Note on an Approximation Treatment for Many-Electron Systems",
		"Phys. Rev.",
		"46",
		"618-622",
		"1934",
		"test tag now"}
)

func TestAddType(t *testing.T) {

	t.Run("Old struct but no type", func(t *testing.T) {
		start := Reference{
			"",
			"MP2",
			"C. M{\\o}ller and M. S. Plesset",
			"Note on an Approximation Treatment for Many-Electron Systems",
			"Phys. Rev.",
			"46",
			"618-622",
			"1934",
			"test tag now"}
		goaltype := "book"
		err := start.AddType(goaltype)
		if err != nil || start[Type] != goaltype {
			t.Errorf("got %s, wanted %s", start[Type], goaltype)
		}
	})

	t.Run("New struct and no type", func(t *testing.T) {
		start := new(Reference)
		goaltype := "book"
		err := start.AddType(goaltype)
		if err != nil || start[Type] != goaltype {
			t.Errorf("got %s, wanted %s", start[Type], goaltype)
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
		"article",
		"MP2",
		"C. M{\\o}ller and M. S. Plesset",
		"Note on an Approximation Treatment for Many-Electron Systems",
		"Phys. Rev.",
		"46",
		"618-622",
		"1934",
		"test tag now"}
	goaltype := "book"
	start.EditType(goaltype)
	if start[Type] != goaltype {
		t.Errorf("got %s, wanted %s", start[Type], goaltype)
	}
}

func TestString(t *testing.T) {
	want := "Type: article\nKey: MP2\nAuthors: C. M{\\o}ller and M. S. Plesset\n" +
		"Title: Note on an Approximation Treatment for Many-Electron Systems\n" +
		"Journal: Phys. Rev.\nVolume: 46\nPages: 618-622\nYear: 1934\nTags: test tag now\n"
	got := ref.String()
	if got != want {
		t.Errorf("\ngot %q\nwad %q", got, want)
	}
}

func TestSearchString(t *testing.T) {
	want := "article, MP2, C. M{\\o}ller and M. S. Plesset, " +
		"Note on an Approximation Treatment for Many-Electron Systems, " +
		"Phys. Rev., 46, 618-622, 1934, test tag now\n"
	got := ref.SearchString()
	if got != want {
		t.Errorf("\ngot %q\nwad %q", got, want)
	}
}

func TestAddRef(t *testing.T) {
	input := ""
	got := AddRef(input)
	want := Reference{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}
