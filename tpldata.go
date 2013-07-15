package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

var (
	ErrEmptyPost          = fmt.Errorf("empty post file")
	ErrInvalidFrontMatter = fmt.Errorf("invalid front matter")
	ErrMissingFrontMatter = fmt.Errorf("missing front matter")
)

// The TemplateData structure contains all the relevant information passed to the
// template to generate the static HTML file.
type TemplateData struct {
	SiteName string
	TagLine  string
	RssURL   string
	Post     *LongPost
	Recent   []*LongPost
	Prev     *ShortPost
	Next     *ShortPost
}

// Create a new TemplateData for the specified post.
func newTemplateData(p *LongPost, i int, r []*LongPost, all []*LongPost) *TemplateData {
	td := &TemplateData{
		SiteName: Options.SiteName,
		TagLine:  Options.TagLine,
		RssURL:   RssURL,
		Post:     p,
		Recent:   r,
	}
	if i > 0 {
		td.Prev = all[i-1].ShortPost
	}
	if i < len(all)-1 {
		td.Next = all[i+1].ShortPost
	}
	return td
}

// The ShortPost structure defines the basic metadata of a post.
type ShortPost struct {
	Slug        string
	Author      string
	Title       string
	Description string
	Lang        string
	PubTime     time.Time
	ModTime     time.Time
}

// The LongPost structure adds the parsed content of the post to the embedded ShortPost information.
type LongPost struct {
	*ShortPost
	Content string
}

// Replace special characters to form a valid slug (post path)
var rxSlug = regexp.MustCompile(`[^a-zA-Z\-_0-9]`)

// Return a valid slug from the file name of the post.
func getSlug(fnm string) string {
	return rxSlug.ReplaceAllString(strings.Replace(fnm, filepath.Ext(fnm), "", 1), "-")
}

// Read the front matter from the post. If there is no front matter, this is
// not a valid post.
func readFrontMatter(s *bufio.Scanner) (map[string]string, error) {
	m := make(map[string]string)
	infm := false
	for s.Scan() {
		l := strings.Trim(s.Text(), " ")
		if l == "---" { // The front matter is delimited by 3 dashes
			if infm {
				// This signals the end of the front matter
				return m, nil
			} else {
				// This is the start of the front matter
				infm = true
			}
		} else if infm {
			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				// Invalid front matter line
				return nil, ErrInvalidFrontMatter
			}
			m[sections[0]] = strings.Trim(sections[1], " ")
		} else if l != "" {
			// No front matter, quit
			return nil, ErrMissingFrontMatter
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return nil, ErrEmptyPost
}

// Create a LongPost from the specified FileInfo.
func newLongPost(fi os.FileInfo) (*LongPost, error) {
	f, err := os.Open(filepath.Join(PostsDir, fi.Name()))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	m, err := readFrontMatter(s)
	if err != nil {
		return nil, err
	}

	slug := getSlug(fi.Name())
	pubdt := fi.ModTime()
	if dt, ok := m["Date"]; ok {
		pubdt, err = time.Parse("2006-01-02", dt)
		if err != nil {
			return nil, err
		}
	}
	sp := &ShortPost{
		slug,
		m["Author"],
		m["Title"],
		m["Description"],
		m["Lang"],
		pubdt,
		fi.ModTime(),
	}

	// Read rest of file
	buf := bytes.NewBuffer(nil)
	for s.Scan() {
		buf.WriteString(s.Text() + "\n")
	}
	if err = s.Err(); err != nil {
		return nil, err
	}
	res := blackfriday.MarkdownCommon(buf.Bytes())
	lp := &LongPost{
		sp,
		string(res),
	}
	return lp, nil
}
