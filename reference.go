package main

import (
	"errors"
	"fmt"
	"regexp"
)

// Errors
var (
	ErrTypeExists    = errors.New("Field Type already exists for this Reference")
	ErrKeyExists     = errors.New("Field Key already exists for this Reference")
	ErrAuthorsExists = errors.New("Field Authors already exists for this Reference")
	ErrTitleExists   = errors.New("Field Title already exists for this Reference")
	ErrJournalExists = errors.New("Field Journal already exists for this Reference")
	ErrVolumeExists  = errors.New("Field Volume already exists for this Reference")
	ErrPagesExists   = errors.New("Field Pages already exists for this Reference")
	ErrYearExists    = errors.New("Field Year already exists for this Reference")
)

// Field enumerates the types of fields a reference can have
type Field int

// Field enumeration
const (
	Type Field = iota
	Key
	Authors
	Title
	Journal
	Volume
	Pages
	Year
	Tags
	NumFields
)

// RefRegex is a combination of a regexp.Regexp and an associated
// Reference Field
type RefRegex struct {
	Expr  *regexp.Regexp
	Value Field
}

// Regular expressions associated with their corresponding Reference
// Fields
var (
	RefRegexes = []RefRegex{
		RefRegex{regexp.MustCompile(`@(article|book|incollection|misc|string){.*`), Type},
		RefRegex{regexp.MustCompile(`(?U)@.*\{(.*),`), Key},
		RefRegex{regexp.MustCompile(`(?iU)Author\s*=\s*{(.*)},`), Authors},
		RefRegex{regexp.MustCompile(`(?iU)Title\s*=\s*{(.*)},`), Title},
		RefRegex{regexp.MustCompile(`(?iU)Journal\s*=\s*{(.*)},`), Journal},
		RefRegex{regexp.MustCompile(`(?i)Volume\s*=\s*\{?([a-z0-9]*)\}?,`), Volume},
		RefRegex{regexp.MustCompile(`(?i)Pages\s*=\s*{?([a-z-0-9]*)}?,`), Pages},
		RefRegex{regexp.MustCompile(`(?i)Year\s*=\s*\{?([a-z0-9]*)}?}`), Year},
		RefRegex{regexp.MustCompile(`(?i)TAGS: (.*\S)`), Tags},
	}
)

// Reference is a wrapper type for a string slice of Fields
type Reference [NumFields]string

// type Reference struct {
// 	Type    string // article, book, etc
// 	Key     string // MP2, Fortenberry07
// 	Authors []string
// 	Title   string
// 	Journal string
// 	Volume  string // for the same as above although less likely
// 	Pages   string // to include the dash and possible alphabetical pages
// 	Year    string
// 	Tags    []string
// }

func (r Reference) String() string {
	return fmt.Sprintf("Type: %s\nKey: %s\nAuthors: %s\nTitle: %s\nJournal: %s\n"+
		"Volume: %s\nPages: %s\nYear: %s\nTags: %s\n", r[Type], r[Key],
		r[Authors], r[Title], r[Journal], r[Volume],
		r[Pages], r[Year], r[Tags])
}

// SearchString formats a reference into a one-line string for use in
// fzf
func (r Reference) SearchString() string {
	s := ""
	for i := range r {
		if r[i] != "" {
			s += fmt.Sprintf("%s", r[i])
			if i < len(r)-1 {
				s += ", "
			}
		}
	}
	s += fmt.Sprintf("\n")
	return s
}

// AddRef adds a reference from the -a command line option
func AddRef(add string) Reference {
	return *new(Reference)
}

// func NewReference(rtype, rkey string, rauth []string, rtitle,
// 	rjour, rvol, rpages, ryear string) Reference {
// 	return Reference{Type: rtype, Key: rkey, Authors: rauth,
// 		Title: rtitle, Journal: rjour, Volume: rvol, Pages: rpages,
// 		Year: ryear}
// }

// AddType adds a reference type to r
func (r *Reference) AddType(refType string) error {
	if r[Type] != "" {
		return ErrTypeExists
	}
	r[Type] = refType
	return nil
}

// EditType allows the type of r to be changed
func (r *Reference) EditType(refType string) {
	r[Type] = refType
}

// func (r *Reference) AddKey(refKey string) error {
// 	if r.Key != "" {
// 		return ErrKeyExists
// 	}
// 	r.Key = refKey
// 	return nil
// }

// func (r *Reference) EditKey(refKey string) {
// 	r.Key = refKey
// }

// func (r *Reference) AddAuthors(refAuthors []string) error {
// 	if !reflect.DeepEqual(r.Authors, []string{""}) {
// 		return ErrAuthorsExists
// 	}
// 	r.Authors = refAuthors
// 	return nil
// }

// func (r *Reference) EditAuthors(refAuthors []string) {
// 	r.Authors = refAuthors
// }

// func (r *Reference) AddTitle(refTitle string) error {
// 	if r.Title != "" {
// 		return ErrTitleExists
// 	}
// 	r.Title = refTitle
// 	return nil
// }

// func (r *Reference) EditTitle(refTitle string) {
// 	r.Title = refTitle
// }

// func (r *Reference) AddJournal(refJournal string) error {
// 	if r.Journal != "" {
// 		return ErrJournalExists
// 	}
// 	r.Journal = refJournal
// 	return nil
// }

// func (r *Reference) EditJournal(refJournal string) {
// 	r.Journal = refJournal
// }

// func (r *Reference) AddVolume(refVolume string) error {
// 	if r.Volume != "" {
// 		return ErrVolumeExists
// 	}
// 	r.Volume = refVolume
// 	return nil
// }

// func (r *Reference) EditVolume(refVolume string) {
// 	r.Volume = refVolume
// }

// func (r *Reference) AddPages(refPages string) error {
// 	if r.Pages != "" {
// 		return ErrPagesExists
// 	}
// 	r.Pages = refPages
// 	return nil
// }

// func (r *Reference) EditPages(refPages string) {
// 	r.Pages = refPages
// }

// func (r *Reference) AddYear(refYear string) error {
// 	if r.Year != "" {
// 		return ErrYearExists
// 	}
// 	r.Year = refYear
// 	return nil
// }

// func (r *Reference) EditYear(refYear string) {
// 	r.Year = refYear
// }

// func (r *Reference) AddTags(refTags []string) error {
// 	r.Tags = append(r.Tags, refTags...)
// 	return nil
// }

// func (r *Reference) EditTags(refTags []string) {
// 	r.Tags = refTags
// }
