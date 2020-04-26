package main

import (
	"reflect"
	"testing"
)

var (
	refs = []Reference{Reference{
		Type:    "article",
		Key:     "MP2",
		Authors: []string{"C. M{\\o}ller", "M. S. Plesset"},
		Title:   "Note on an Approximation Treatment for Many-Electron Systems",
		Journal: "Phys. Rev.",
		Volume:  "46",
		Pages:   "618-622",
		Year:    "1934",
		Tags:    []string{"test", "tag", "now"}}}
)

func TestReadBib(t *testing.T) {
	got := ReadBib("tex/min.bib")
	if !reflect.DeepEqual(got, refs) {
		t.Errorf("got %#v, wanted %#v\n", got, refs)
	}
}

func TestMakeBib(t *testing.T) {
	got := MakeBib(refs)
	want := []string{"@article{MP2,",
		"Author={C. M{\\o}ller and M. S. Plesset},",
		"Title={Note on an Approximation Treatment for Many-Electron Systems},",
		"Journal={Phys. Rev.},",
		"Volume=46,",
		"Pages={618-622},",
		"Year=1934}",
		"TAGS: test tag now",
		""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v\n", got, want)
	}
}

func TestWriteBib(t *testing.T) {
	WriteBib(refs, "tex/testbib.out")
}
