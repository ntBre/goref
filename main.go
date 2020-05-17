package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Read these from config file
var (
	Config [noptions]string
)

/*
Reference could be an array instead of a struct
type Field int
const (
Type Field = iota
Key
Author
...
Tags
NumFields // nice way to get the number of fields, just leave at end
)
func (f Field) String() string {} // for printing 
This would require all of the fields to be stored as strings
which is fine, just split Authors on ` and ` if needed and Tags on space
Reference[Type] = Field
type Field struct {
    *regexp.Regexp
    Value
}
but then how to associate the regular expressions with the right positions?
just have to maintain an array of them
parallel arrays of Reference and regexes and loop through regexes, use same
index in the Reference array
    - negates need for this Field struct
at every refstring, loop through regexes and Reference[i] = regex[i].FindMatch
refs is still []Reference, so append at end
Field alias above makes this not work for integer i without a cast, so maybe no alias
go back to custom Regex type like in go-cart infile reading
Regex and Field value, then the array of regexes doesnt have to be ordered either

The way I've done it already makes the most sense to me intuitively, but I dislike
the huge tree of ifs when theyre all doing the same thing
Similarly the list of named regular expressions; seems clear that there should
be some way to pair the regexes and the Reference fields
*/

func ReadBib(bibname string) (refs []Reference) {

	var noTags bool
	file, err := ioutil.ReadFile(bibname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")

	nref := -1
	reftype := regexp.MustCompile(`@(article|book|incollection|misc|string){.*`)
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
	// find where references start
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
		if strings.Contains(line, "@") {
			indices = append(indices, i)
		}
	}
	// put references on single lines
	refstrings := make([]string, 0)
	for i, _ := range indices {
		if i < len(indices)-1 {
			refstrings = append(refstrings, strings.Join(lines[indices[i]:indices[i+1]], " "))
		} else {
			refstrings = append(refstrings, strings.Join(lines[indices[i]:], " "))
		}
	}

	for _, ref := range refstrings {
		if reftype.MatchString(ref) {
			refs = append(refs, *new(Reference))
			nref++
			noTags = true
			refs[nref].Type = reftype.FindStringSubmatch(ref)[1]
		}
		if key.MatchString(ref) {
			// TODO farm this to helper with error handling
			// grabbing first match of substring
			refs[nref].Key = key.FindStringSubmatch(ref)[1]
		}
		if author.MatchString(ref) {
			refs[nref].Authors = and.Split(author.FindStringSubmatch(ref)[1], -1)
		}
		if title.MatchString(ref) {
			refs[nref].Title = title.FindStringSubmatch(ref)[1]
		}
		if journal.MatchString(ref) {
			refs[nref].Journal = journal.FindStringSubmatch(ref)[1]
		}
		if pages.MatchString(ref) {
			refs[nref].Pages = pages.FindStringSubmatch(ref)[1]
		}
		if volume.MatchString(ref) {
			refs[nref].Volume = volume.FindStringSubmatch(ref)[1]
		}
		if year.MatchString(ref) {
			refs[nref].Year = year.FindStringSubmatch(ref)[1]
		}
		if tags.MatchString(ref) {
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
			fmt.Sprintf("TAGS: %s", strings.Join(ref.Tags, " ")), "")
	}
	return
}

func WriteBib(refs []Reference, filename string) {
	lines := strings.Join(MakeBib(refs), "\n")
	ioutil.WriteFile(filename, []byte(lines), 0755)
}

func WriteFZFList(refs []Reference, w io.Writer) {
	lines := ""
	for _, ref := range refs {
		lines += ref.SearchString()
	}
	w.Write([]byte(lines))
}

func FuzzyFind(refs []Reference) string {
	// fuzzy find in refs and return the found string
	buf := new(bytes.Buffer)
	WriteFZFList(refs, buf)
	command := "fzf -m"
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	// unclear to me why go routine necessary but it seems to be
	go func() {
		fmt.Fprintln(in, buf.String())
		in.Close()
	}()
	result, _ := cmd.Output()
	return strings.TrimSpace(string(result))
	
}

func main() {
	// for building executable to use now
	// ouch, going to have to handle different systems' directory structures
	ParseConfig("/home/brent/.config/goref/config")
	refs := ReadBib(Config[bibfile])
	fmt.Println(FuzzyFind(refs))
}

// func main() {
// 	box := tview.NewBox().SetBorder(true).SetTitle("[blue::l]Hello world!!")
// 	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
// 		panic(err)
// 	}
// }
