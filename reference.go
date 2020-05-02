package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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

type Reference struct {
	Type    string // article, book, etc
	Key     string // MP2, Fortenberry07
	Authors []string
	Title   string
	Journal string
	Volume  string // for the same as above although less likely
	Pages   string // to include the dash and possible alphabetical pages
	Year    string
	Tags    []string
}

func (r Reference) String() string {
	return fmt.Sprintf("Type: %s\nKey: %s\nAuthors: %s\nTitle: %s\nJournal: %s\n"+
		"Volume: %s\nPages: %s\nYear: %s\nTags: %s\n", r.Type, r.Key,
		strings.Join(r.Authors, " and "), r.Title, r.Journal, r.Volume,
		r.Pages, r.Year, strings.Join(r.Tags, ", "))
}

func (r Reference) SearchString() string {
	return fmt.Sprintf("%s: %s, %s, %s %s, %s, %s; %s\n",
		r.Key, strings.Join(r.Authors, " and "),
		r.Title, r.Journal, r.Volume,
		r.Pages, r.Year, strings.Join(r.Tags, ", "))
}

func NewReference(rtype, rkey string, rauth []string, rtitle,
	rjour, rvol, rpages, ryear string) Reference {
	return Reference{Type: rtype, Key: rkey, Authors: rauth,
		Title: rtitle, Journal: rjour, Volume: rvol, Pages: rpages,
		Year: ryear}
}

func (r *Reference) AddType(refType string) error {
	if r.Type != "" {
		return ErrTypeExists
	}
	r.Type = refType
	return nil
}

func (r *Reference) EditType(refType string) {
	r.Type = refType
}

func (r *Reference) AddKey(refKey string) error {
	if r.Key != "" {
		return ErrKeyExists
	}
	r.Key = refKey
	return nil
}

func (r *Reference) EditKey(refKey string) {
	r.Key = refKey
}

func (r *Reference) AddAuthors(refAuthors []string) error {
	if !reflect.DeepEqual(r.Authors, []string{""}) {
		return ErrAuthorsExists
	}
	r.Authors = refAuthors
	return nil
}

func (r *Reference) EditAuthors(refAuthors []string) {
	r.Authors = refAuthors
}

func (r *Reference) AddTitle(refTitle string) error {
	if r.Title != "" {
		return ErrTitleExists
	}
	r.Title = refTitle
	return nil
}

func (r *Reference) EditTitle(refTitle string) {
	r.Title = refTitle
}

func (r *Reference) AddJournal(refJournal string) error {
	if r.Journal != "" {
		return ErrJournalExists
	}
	r.Journal = refJournal
	return nil
}

func (r *Reference) EditJournal(refJournal string) {
	r.Journal = refJournal
}

func (r *Reference) AddVolume(refVolume string) error {
	if r.Volume != "" {
		return ErrVolumeExists
	}
	r.Volume = refVolume
	return nil
}

func (r *Reference) EditVolume(refVolume string) {
	r.Volume = refVolume
}

func (r *Reference) AddPages(refPages string) error {
	if r.Pages != "" {
		return ErrPagesExists
	}
	r.Pages = refPages
	return nil
}

func (r *Reference) EditPages(refPages string) {
	r.Pages = refPages
}

func (r *Reference) AddYear(refYear string) error {
	if r.Year != "" {
		return ErrYearExists
	}
	r.Year = refYear
	return nil
}

func (r *Reference) EditYear(refYear string) {
	r.Year = refYear
}

func (r *Reference) AddTags(refTags []string) error {
	r.Tags = append(r.Tags, refTags...)
	return nil
}

func (r *Reference) EditTags(refTags []string) {
	r.Tags = refTags
}
