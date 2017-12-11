package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type cli struct {
	token   string
	query   string
	dest    string
	extract *regexp.Regexp // Maybe nil
}

func (c *cli) ensureReposDir() error {
	s, err := os.Stat(c.dest)
	if err != nil {
		return os.Mkdir(c.dest, 0755)
	}
	if !s.IsDir() {
		return fmt.Errorf("Cannot create directory '%s' because it's a file", c.dest)
	}
	return nil
}

func (c *cli) run() (err error) {
	if err = c.ensureReposDir(); err != nil {
		return
	}
	col := newCollector(c.query, c.token, c.dest, c.extract, nil)
	_, _, err = col.collect()
	return
}

func newCLI(t, q, l, d, e string) (*cli, error) {
	var err error

	if env := os.Getenv("GITHUB_TOKEN"); env != "" && t == "" {
		t = env
	}

	if t == "" || l == "" {
		return nil, fmt.Errorf("API token and language must be set. Please see -help for more detail")
	}

	if d == "" {
		d, err = os.Getwd()
		if err != nil {
			return nil, err
		}
		d = filepath.Join(d, "repos")
	}

	var r *regexp.Regexp
	if e != "" {
		r, err = regexp.Compile(e)
		if err != nil {
			return nil, err
		}
	}

	q = fmt.Sprintf("%s language:%s fork:false", q, l)
	return &cli{t, q, d, r}, nil
}
