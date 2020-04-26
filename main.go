package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Reference struct {
	Type    string // @article, book, etc
	Key     string // MP2, Fortenberry07
	Authors []string
	Title   string
	Journal string
	Pages   string // to include the dash and possible alphabetical pages
	Volume  string // for the same as above although less likely
	Year    string
	Tags    []string
}

func ReplaceSubex(re *regexp.Regexp, s string, n int) string {
	// return the nth subexpression re of s
	return strings.TrimSpace(string(re.ReplaceAllString(s, "$"+strconv.Itoa(n))))
}

func SplitAndTrim(s string, re *regexp.Regexp) []string {
	return re.Split(strings.TrimSpace(s), -1)
}

func ReadBib(bibname string) (refs []Reference) {

	// TODO handle file with fields over multiple lines
	file, err := ioutil.ReadFile(bibname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")

	nref := 0
	reftype := regexp.MustCompile(`@(article|book){.*`)
	key := regexp.MustCompile(`@.*{(MP2),`)
	// TODO make this more general, reliant on case now
	author := regexp.MustCompile(`Author={(.*)},`)
	and := regexp.MustCompile(`\s+and\s+`)
	title := regexp.MustCompile(`Title={(.*)},`)
	journal := regexp.MustCompile(`Journal={(.*)},`)
	volume := regexp.MustCompile(`Volume=(.*),`)
	pages := regexp.MustCompile(`Pages={(.*)},`)
	// TODO specifically here requires brace on end
	year := regexp.MustCompile(`Year=(.*)}`)
	tags := regexp.MustCompile(`TAGS: (.*)`)
	tagspace := regexp.MustCompile(` `)

	for _, line := range lines {
		noTags := true
		switch {
		case reftype.MatchString(line):
			refs = append(refs, *new(Reference))
			refs[nref].Type = ReplaceSubex(reftype, line, 1)
			fallthrough
		case key.MatchString(line):
			refs[nref].Key = ReplaceSubex(key, line, 1)
		case author.MatchString(line):
			refs[nref].Authors = SplitAndTrim(ReplaceSubex(author, line, 1), and)
		case title.MatchString(line):
			refs[nref].Title = ReplaceSubex(title, line, 1)
		case journal.MatchString(line):
			refs[nref].Journal = ReplaceSubex(journal, line, 1)
		case pages.MatchString(line):
			refs[nref].Pages = ReplaceSubex(pages, line, 1)
		case volume.MatchString(line):
			refs[nref].Volume = ReplaceSubex(volume, line, 1)
		case year.MatchString(line):
			refs[nref].Year = ReplaceSubex(year, line, 1)
		case tags.MatchString(line):
			refs[nref].Tags = SplitAndTrim(ReplaceSubex(tags, line, 1), tagspace)
			noTags = false
		}
		if noTags {
			refs[nref].Tags = []string{""}
		}
	}
	return
}

func MakeBib(refs []Reference) (lines []string) {
	for _, ref := range refs {
		// TODO try for := range reflect.Fields, not sure if it will keep order
		lines = append(lines, fmt.Sprintf("@%s{%s,", ref.Type, ref.Key),
			fmt.Sprintf("Author={%s},", strings.Join(ref.Authors, " and ")),
			fmt.Sprintf("Title={%s},", ref.Title),
			fmt.Sprintf("Journal={%s},", ref.Journal),
			fmt.Sprintf("Volume=%s,", ref.Volume),
			fmt.Sprintf("Pages={%s},", ref.Pages),
			fmt.Sprintf("Year=%s}", ref.Year),
			fmt.Sprintf("TAGS: %s", strings.Join(ref.Tags, " ")),
			"")
	}
	return
}

func WriteBib(refs []Reference, filename string) {
	lines := strings.Join(MakeBib(refs), "\n")
	ioutil.WriteFile(filename, []byte(lines), 0755)
}
