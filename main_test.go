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
	refs2 = []Reference{
		Reference{
			Type:    "article",
			Key:     "MP2",
			Authors: []string{"C. M{\\o}ller", "M. S. Plesset"},
			Title:   "Note on an Approximation Treatment for Many-Electron Systems",
			Journal: "Phys. Rev.",
			Volume:  "46",
			Pages:   "618-622",
			Year:    "1934",
			Tags:    []string{"test", "tag", "now"},
		},
		Reference{
			Type:    "article",
			Key:     "Rittby91",
			Authors: []string{"C. M. L. Rittby"},
			Title:   "",
			Journal: "J. Chem. Phys.",
			Volume:  "95",
			Pages:   "5609-5611",
			Year:    "1991",
			Tags:    []string{""},
		},
	}
	refpg = []Reference{
		Reference{
			Type:    "article",
			Key:     "Huang08",
			Authors: []string{"X. Huang", "T. J. Lee"},
			Title:   "A Procedure for Computing Accurate $Ab\\ Initio$ "+
				"Quartic Force Fields: Application to {HO$_2$$^+$ and H$_2$O}",
			Journal: "J. Chem. Phys.",
			Volume:  "129",
			Pages:   "044312",
			Year:    "2008",
			Tags:    []string{""},
			
		},
	}
)

func TestReadBib(t *testing.T) {
	t.Run("basic reference", func(t *testing.T) {
		got := ReadBib("tex/min.bib")
		if !reflect.DeepEqual(got, refs) {
			t.Errorf("\ngot %q\nwad %q\n", got, refs)
		}
	})
	t.Run("two references", func(t *testing.T) {
		got := ReadBib("tex/med.bib")
		if !reflect.DeepEqual(got, refs2) {
			t.Errorf("\ngot %q\nwad %q\n", got, refs2)
		}
	})
	t.Run("one page number", func(t *testing.T) {
		got := ReadBib("tex/onepg.bib")
		if !reflect.DeepEqual(got, refpg) {
			t.Errorf("\ngot %q\nwad %q\n", got, refpg)
		}
	})
	// TODO continue looking for problems with this
	// fmt.Printf("%s", ReadBib("tex/refs.bib"))
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
