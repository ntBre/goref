package main

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

var (
	refs = []Reference{
		Reference{"article",
			"MP2",
			"C. M{\\o}ller and M. S. Plesset",
			"Note on an Approximation Treatment for Many-Electron Systems",
			"Phys. Rev.",
			"46",
			"618-622",
			"1934",
			"test tag now"},
	}
	refs2 = []Reference{
		Reference{"article",
			"MP2",
			"C. M{\\o}ller and M. S. Plesset",
			"Note on an Approximation Treatment for Many-Electron Systems",
			"Phys. Rev.",
			"46",
			"618-622",
			"1934",
			"test tag now"},
		Reference{
			"article",
			"Rittby91",
			"C. M. L. Rittby",
			"",
			"J. Chem. Phys.",
			"95",
			"5609-5611",
			"1991",
			"",
		},
	}
	refpg = []Reference{
		Reference{
			"article",
			"Huang08",
			"X. Huang and T. J. Lee",
			"A Procedure for Computing Accurate $Ab\\ Initio$ " +
				"Quartic Force Fields: Application to {HO$_2$$^+$ and H$_2$O}",
			"J. Chem. Phys.",
			"129",
			"044312",
			"2008",
			"",
		},
	}
	book = []Reference{
		Reference{
			"book",
			"Cook",
			"D. B. Cook",
			"Handbook of Computational Quantum Chemistry",
			"",
			"",
			"",
			"2005",
			"",
		},
		Reference{
			"book",
			"Cramer",
			"C. J. Cramer",
			"Essentials of Computational Chemistry: Theories and Models",
			"",
			"",
			"",
			"2004",
			"",
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
	// remember to make tags empty string, different from nil
	t.Run("two books", func(t *testing.T) {
		got := ReadBib("tex/book.bib")
		if !reflect.DeepEqual(got, book) {
			t.Errorf("\ngot %q\nwad %q\n", got, book)
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
		t.Errorf("\ngot %#v, \nwad %#v\n", got, want)
	}
}

func TestWriteBib(t *testing.T) {
	WriteBib(refs, "tex/testbib.out")
}

func TestWriteFZFList(t *testing.T) {
	refs := ReadBib("tex/refs.bib")
	f, _ := os.Create("tex/test.fzf")
	WriteFZFList(refs, bufio.NewWriter(f))
}
