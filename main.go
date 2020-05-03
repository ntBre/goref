package main

import (
	"fmt"
	"github.com/rivo/tview"
	"io/ioutil"
	"regexp"
	"strings"
)

func ReadBib(bibname string) (refs []Reference) {

	var noTags bool
	file, err := ioutil.ReadFile(bibname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")

	nref := -1
	reftype := regexp.MustCompile(`@(article|book){.*`)
	key := regexp.MustCompile(`(?U)@.*\{(.*),`)
	author := regexp.MustCompile(`(?iU)Author\s*=\s*{(.*)},`)
	and := regexp.MustCompile(`\s+and\s+`)
	title := regexp.MustCompile(`(?iU)Title\s*=\s*{(.*)},`)
	journal := regexp.MustCompile(`(?iU)Journal\s*=\s*{(.*)},`)
	volume := regexp.MustCompile(`(?i)Volume\s*=\s*\{?([a-z0-9]*)\}?,`)
	pages := regexp.MustCompile(`(?i)Pages\s*=\s*{?([a-z-0-9]*)}?,`)
	year := regexp.MustCompile(`(?i)Year\s*=\s*\{?([a-z0-9]*)}?}`)
	// end with nonspace
	tags := regexp.MustCompile(`(?i)TAGS: (.*\S)`)
	tagspace := regexp.MustCompile(` `)
	// eqBracket := regexp.MustCompile(`\s*=\s*{\s*`)
	indices := make([]int, 0)
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
		if strings.Contains(line, "@") {
			indices = append(indices, i)
		}
	}
	refstrings := make([]string, 0)
	for i, _ := range indices {
		if i < len(indices)-1 {
			refstrings = append(refstrings, strings.Join(lines[indices[i]:indices[i+1]], " "))
		} else {
			refstrings = append(refstrings, strings.Join(lines[indices[i]:], " "))
		}
	}

	for _, ref := range refstrings {
		if reftype.MatchString(ref){
			refs = append(refs, *new(Reference))
			nref++
			noTags = true
			refs[nref].Type = reftype.FindStringSubmatch(ref)[1]
		}
		if key.MatchString(ref){
			refs[nref].Key = key.FindStringSubmatch(ref)[1]
		}
		if author.MatchString(ref) {
			refs[nref].Authors = and.Split(author.FindStringSubmatch(ref)[1], -1)
		}
		if title.MatchString(ref) {
			refs[nref].Title = title.FindStringSubmatch(ref)[1]
		}
		if journal.MatchString(ref){
			refs[nref].Journal = journal.FindStringSubmatch(ref)[1]
		}
		if pages.MatchString(ref){
			refs[nref].Pages = pages.FindStringSubmatch(ref)[1]
		}
		if volume.MatchString(ref){
			refs[nref].Volume = volume.FindStringSubmatch(ref)[1]
		}
		if year.MatchString(ref){
			refs[nref].Year = year.FindStringSubmatch(ref)[1]
		}
		if tags.MatchString(ref){
			refs[nref].Tags = tagspace.Split(tags.FindStringSubmatch(ref)[1], -1)
			noTags = false
		}
		if noTags {
			refs[nref].Tags = []string{""}
		}
	}
	return
}

func MakeBib(refs []Reference) (lines []string) {
	// TODO can probably refactor this with reference.String()
	// well something similar
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

func WriteFZFList(refs []Reference, filename string) {
	lines := ""
	for _, ref := range refs {
		lines += ref.SearchString()
	}
	ioutil.WriteFile(filename, []byte(lines), 0755)
}

func main() {
	box := tview.NewBox().SetBorder(true).SetTitle("[blue::l]Hello world!!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
