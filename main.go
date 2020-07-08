package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Command line flags
var (
	search = flag.Bool("s", false, "search the bib file with fzf and print the match")
	open   = flag.Bool("o", false, "like -s but try to open the chosen file in pdfviewer")
	add    = flag.String("a", "", "add a reference in \"Authors, Title, Journal, Volume, Page, Year\" format")
)

// Global configuration and reference arrays
var (
	Config [noptions]string
	refs   []Reference
)

// ReadBib reads a LaTeX bibliography file into a slice of References
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
	for i := range indices {
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

// MakeBib converts a slice of Reference back into the lines of a
// LaTeX bibliography file
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

// WriteBib uses MakeBib to write a LaTeX bibliography file
func WriteBib(refs []Reference, filename string) {
	lines := strings.Join(MakeBib(refs), "\n")
	ioutil.WriteFile(filename, []byte(lines), 0755)
}

// WriteFZFList writes a slice of Reference into a format useable by
// fzf
func WriteFZFList(refs []Reference, w io.Writer) {
	lines := ""
	for _, ref := range refs {
		lines += ref.SearchString()
	}
	w.Write([]byte(lines))
}

// FuzzyFind is the fzf interface
func FuzzyFind(refs []Reference) string {
	// fuzzy find in refs and return the found string
	buf := new(bytes.Buffer)
	WriteFZFList(refs, buf)
	// search from the top
	// (e)xact search, search for words not letters
	command := "fzf -m --layout=reverse -e"
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

// ParseFlags parses command line flags and returns the remaining
// arguments
func ParseFlags() {
	flag.Parse()
}

// CleanUp writes the updated LaTeX bibliography file and exits with a
// status corresponding to err
func CleanUp(err error) {
	var exitCode int
	if err != nil {
		fmt.Fprintf(os.Stderr, "goref: %s\n", err)
		exitCode = 2
	}
	WriteBib(refs, Config[bibfile])
	os.Exit(exitCode)
}

func main() {
	ParseFlags()
	ParseConfig("/home/brent/.config/goref/config")
	refs = ReadBib(Config[bibfile])
	// restructure inside here as function calls
	// so they can be called by TUI buttons too
	switch {
	case *search || *open:
		print := true
		out := FuzzyFind(refs)
		if *open {
			match := regexp.MustCompile(`(?i)^[a-z]+, ([a-z0-9]+)`).
				FindStringSubmatch(out)
			var pdfFile string
			if len(match) >= 2 {
				pdfFile = match[1]
			}
			pdfPath := Config[library] + "/" + pdfFile + ".pdf"
			if _, err := os.Stat(pdfPath); !os.IsNotExist(err) {
				exec.Command(Config[pdfcmd], pdfPath).Start()
				print = false
			}
		}
		if print && out != "" {
			fmt.Println(out)
		}
	case *add != "":
		// AddRef(*add)
		fmt.Println(*add)
	default:
		tview.DoubleClickInterval = 0
		app := tview.NewApplication().EnableMouse(true)
		text := tview.NewTextView().
			SetChangedFunc(func() {
				app.Draw()
			})
		text.SetBorder(true).SetTitle("References")
		for _, ref := range refs {
			fmt.Fprintf(text, ref.SearchString())
		}
		flex := tview.NewFlex().
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Menu"), 20, 1, false).
			AddItem(text, 0, 1, true)
		// have to duplicate this for every box, or put it on the app,
		// which takes over control of the subelement keybinds
		text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'q':
				app.Stop()
				return nil
			}
			return nil
		})
		if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
			CleanUp(err)
		}
	}
	CleanUp(nil)
}
