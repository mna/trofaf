package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type TemplateData struct {
	SiteName string
	Post     *LongPost
	Recent   []*ShortPost
	Prev     *ShortPost
	Next     *ShortPost
}

func newTemplateData(p *LongPost, i int, r []*ShortPost, all []*LongPost) *TemplateData {
	td := &TemplateData{SiteName: SiteName, Post: p, Recent: r}

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

func newLongPost(fi os.FileInfo) *LongPost {
	slug := getSlug(fi.Name())
	sp := &ShortPost{
		slug,
		"author",      // TODO : Complete...
		slug,          // TODO : Read first heading, or front matter
		"description", // TODO : Read front matter
		fi.ModTime(),  // TODO : This is NOT the pub time...
		fi.ModTime(),
	}

	f, err := os.Open(filepath.Join(PostsDir, fi.Name()))
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	res := blackfriday.MarkdownCommon(b)
	lp := &LongPost{
		sp,
		string(res),
	}
	return lp
}

func (lp *LongPost) Short() *ShortPost {
	return lp.ShortPost
}
