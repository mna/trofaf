package main

import (
	"bufio"
	"bytes"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type TemplateData struct {
	SiteName string
	TagLine  string
	RssURL   string
	Post     *LongPost
	Recent   []*ShortPost
	Prev     *ShortPost
	Next     *ShortPost
}

func newTemplateData(p *LongPost, i int, r []*ShortPost, all []*LongPost) *TemplateData {
	b, err := url.Parse(Options.BaseURL)
	if err != nil {
		panic(err) // TODO : Manage errors
	}
	u, err := b.Parse("/rss")
	if err != nil {
		panic(err)
	}
	td := &TemplateData{
		SiteName: Options.SiteName,
		TagLine:  Options.TagLine,
		RssURL:   u.String(),
		Post:     p,
		Recent:   r,
	}
	if i > 0 {
		td.Prev = all[i-1].Short()
	}
	if i < len(all)-2 {
		td.Next = all[i+1].Short()
	}
	return td
}

type ShortPost struct {
	Slug        string
	Author      string
	Title       string
	Description string
	Lang        string
	PubTime     time.Time
	ModTime     time.Time
}

type LongPost struct {
	*ShortPost
	Content string
}

var rxSlug = regexp.MustCompile(`[^a-zA-Z\-_0-9]`)

func getSlug(fnm string) string {
	return rxSlug.ReplaceAllString(strings.Replace(fnm, filepath.Ext(fnm), "", 1), "-")
}

func readFrontMatter(s *bufio.Scanner) map[string]string {
	m := make(map[string]string)
	infm := false
	for s.Scan() {
		l := strings.Trim(s.Text(), " ")
		if l == "---" { // The front matter is delimited by 3 dashes
			if infm {
				// This signals the end of the front matter
				return m
			} else {
				// This is the start of the front matter
				infm = true
			}
		} else if infm {
			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				// Invalid front matter line
				log.Println("POST ERROR invalid front matter line: ", l)
				return nil
			}
			m[sections[0]] = strings.Trim(sections[1], " ")
		} else if l != "" {
			// No front matter, quit
			return nil
		}
	}
	if infm {
		log.Println("POST ERROR unclosed front matter")
	} else if err := s.Err(); err != nil {
		log.Println("POST ERROR ", err)
	}
	return nil
}

func newLongPost(fi os.FileInfo) *LongPost {
	log.Println("processing post ", fi.Name())
	f, err := os.Open(filepath.Join(PostsDir, fi.Name()))
	if err != nil {
		log.Println("POST ERROR ", err)
		return nil
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	m := readFrontMatter(s)

	slug := getSlug(fi.Name())
	pubdt := fi.ModTime()
	if dt, ok := m["Date"]; ok {
		pubdt, err = time.Parse("2006-01-02", dt)
		if err != nil {
			log.Println("POST ERROR invalid date format: ", dt)
			pubdt = fi.ModTime()
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
		log.Println("POST ERROR ", err)
		return nil
	}
	res := blackfriday.MarkdownCommon(buf.Bytes())
	lp := &LongPost{
		sp,
		string(res),
	}
	return lp
}

func (lp *LongPost) Short() *ShortPost {
	return lp.ShortPost
}
