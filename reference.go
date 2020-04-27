package main

import (
	"errors"
)

var (
	ErrTypeExists = errors.New("Field Type already exists for this Reference")
)

type Reference struct {
	Type    string // article, book, etc
	Key     string // MP2, Fortenberry07
	Authors []string
	Title   string
	Journal string
	Pages   string // to include the dash and possible alphabetical pages
	Volume  string // for the same as above although less likely
	Year    string
	Tags    []string
}

func (r *Reference) AddType(reftype string) error {
	if r.Type != "" {
		return ErrTypeExists
	}
	r.Type = reftype
	return nil
}

func (r *Reference) EditType(reftype string) {
	r.Type = reftype
}
