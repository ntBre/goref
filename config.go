package main

import (
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	library Option = iota
	bibfile
	pdfcmd
	noptions
)

type Option int

func (o Option) String() string {
	return []string{
		"library",
		"bibfile",
		"pdfcmd",
	}[o]
}

type ConfigRegex struct {
	Expr *regexp.Regexp
	Opt  Option
}

var (
	Regexes = []ConfigRegex{
		ConfigRegex{regexp.MustCompile(`(?i)library=`), library},
		ConfigRegex{regexp.MustCompile(`(?i)refs=`), bibfile},
		ConfigRegex{regexp.MustCompile(`(?i)pdfcmd=`), pdfcmd},
	}
)

func ParseConfig(config string) {
	file, err := ioutil.ReadFile(config)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		for _, regex := range Regexes {
			if regex.Expr.MatchString(line) {
				split := strings.Split(line, "=")
				Config[regex.Opt] = split[len(split)-1]
			}
		}
	}
}
