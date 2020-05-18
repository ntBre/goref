package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/rivo/tview"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	Config [noptions]string
	search bool
	open   bool
	help   bool
)

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

type Field int

type RefRegex struct {
	Expr  *regexp.Regexp
	Value Field
}

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

func ReadBib(bibname string) (refs []Reference) {

	file, err := ioutil.ReadFile(bibname)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")

	nref := -1

	// and := regexp.MustCompile(`\s+and\s+`)
	// end with nonspace
	// tagspace := regexp.MustCompile(` `)
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
		for _, regex := range RefRegexes {
			if regex.Expr.MatchString(ref) {
				if regex.Value == Type {
					refs = append(refs, *new(Reference))
					nref++
				}
				// no longer splitting tags and authors, all strings
				// no longer need noTags, let default tags in array be ""
				//     - automatic from the initialization of []string
				refs[nref][regex.Value] = regex.Expr.FindStringSubmatch(ref)[1]
			}
		}
	}
	return
}

func MakeBib(refs []Reference) (lines []string) {
	for _, ref := range refs {
		lines = append(lines, fmt.Sprintf("@%s{%s,", ref[Type], ref[Key]),
			fmt.Sprintf("Author={%s},", ref[Authors]),
			fmt.Sprintf("Title={%s},", ref[Title]),
			fmt.Sprintf("Journal={%s},", ref[Journal]),
			fmt.Sprintf("Volume=%s,", ref[Volume]),
			fmt.Sprintf("Pages={%s},", ref[Pages]),
			fmt.Sprintf("Year=%s}", ref[Year]),
			fmt.Sprintf("TAGS: %s", ref[Tags]), "")
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
	// search from the top
	command := "fzf -m --layout=reverse"
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

func ParseFlag() {
	flag.BoolVar(&search, "s", false, "instead of running the application"+
		" just search the bib file with fzf")
	flag.BoolVar(&open, "o", false, "with -s, try to open the chosen file in pdfviewer")
	flag.BoolVar(&help, "h", false, "list the command line options")
	flag.Parse()
}

func main() {
	ParseFlag()
	ParseConfig("/home/brent/.config/goref/config")
	refs := ReadBib(Config[bibfile])
	if search {
		out := FuzzyFind(refs)
		// need a way to grab Key and add .pdf
		if open {
			pdfFile := regexp.MustCompile(`(?i)^[a-z]+, ([a-z0-9]+)`).FindStringSubmatch(out)[1]
			pdfPath := Config[library] + "/" + pdfFile + ".pdf"
			if _, err := os.Stat(pdfPath); !os.IsNotExist(err) {
				exec.Command(Config[pdfcmd], pdfPath).Run()
			} else {
				if out != "" {
					fmt.Println(out)
				}
			}
		}
	} else if help {
		flag.PrintDefaults()
	} else {
		box := tview.NewBox().SetBorder(true).SetTitle("[blue]goref")
		if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
			panic(err)
		}
	}
	WriteBib(refs, Config[bibfile])
}
