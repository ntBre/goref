package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func ReplaceSubex(re *regexp.Regexp, s string, n int) string {
	// return the nth subexpression re of s
	return strings.TrimSpace(string(re.ReplaceAllString(s, "$"+strconv.Itoa(n))))
}

func SplitAndTrim(s string, re *regexp.Regexp) []string {
	return re.Split(strings.TrimSpace(s), -1)
}

func ReadBib(bibname string) (refs []Reference) {

	var noTags bool
	file, err := ioutil.ReadFile(bibname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")

	nref := -1
	reftype := regexp.MustCompile(`@(article|book){.*`)
	key := regexp.MustCompile(`@.*{(.*),`)
	author := regexp.MustCompile(`(?i)Author\s*=\s*{(.*)},`)
	and := regexp.MustCompile(`\s+and\s+`)
	title := regexp.MustCompile(`(?i)Title\s*=\s*{(.*)},`)
	journal := regexp.MustCompile(`(?i)Journal\s*=\s*{(.*)},`)
	volume := regexp.MustCompile(`(?i)Volume\s*=\s*\{?([a-z0-9]*)\}?,`)
	pages := regexp.MustCompile(`(?i)Pages\s*=\s*{?([a-z-0-9]*)}?,`)
	year := regexp.MustCompile(`(?i)Year\s*=\s*\{?([a-z0-9]*)}?}`)
	tags := regexp.MustCompile(`(?i)TAGS: (.*)`)
	tagspace := regexp.MustCompile(` `)
	eqBracket := regexp.MustCompile(`\s*=\s*{\s*`)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if eqBracket.MatchString(line) && !strings.Contains(line, "},") {
			for !strings.Contains(line, "},") {
				line += " " + strings.TrimSpace(lines[i+1])
				i++
			}
		}
		switch {
		case reftype.MatchString(line):
			refs = append(refs, *new(Reference))
			nref++
			noTags = true
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
