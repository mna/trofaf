package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/russross/blackfriday"
)

const (
	// Sometimes many events can be triggered in succession for the same file
	// (i.e. Create followed by Modify, etc.). No need to rush to generate
	// the HTML, just wait for it to calm down before processing.
	watchEventDelay = 30 * time.Second
	maxRecentPosts  = 2
)

// Receive watcher events for the posts directory. All events require re-generating
// the whole site (because the template may display the n most recent posts, the
// next and previous post, etc.). It could be fine-tuned based on what data we give
// to the templates, but for now, lazy approach.
func watch(w *fsnotify.Watcher) {
	var delay <-chan time.Time
	for {
		select {
		case ev := <-w.Event:
			log.Print("watch event ", ev)
			// Regenerate the files after the delay, reset the delay if an event is triggered
			// in the meantime
			delay = time.After(watchEventDelay)

		case err := <-w.Error:
			log.Print("WATCH ERROR ", err)

		case <-delay:
			log.Print("trigger regeneration of site")
			regeneratePosts()
		}
	}
}

type sortableFileInfo []os.FileInfo

func (s sortableFileInfo) Len() int           { return len(s) }
func (s sortableFileInfo) Less(i, j int) bool { return s[i].ModTime().Before(s[j].ModTime()) }
func (s sortableFileInfo) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func FilterDir(s sortableFileInfo) sortableFileInfo {
	for i := 0; i < len(s); {
		if s[i].IsDir() {
			s[i], s = s[len(s)-1], s[:len(s)-1]
		} else {
			i++
		}
	}
	return s
}

func regeneratePosts() {
	// Clear the public directory, except subdirs
	fis, err := ioutil.ReadDir(PublicDir)
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	for _, fi := range fis {
		if !fi.IsDir() {
			err = os.Remove(filepath.Join(PublicDir, fi.Name()))
			if err != nil {
				log.Println("DELETE ERROR ", err)
			}
		}
	}
	// Now read the posts
	fis, err = ioutil.ReadDir(PostsDir)
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	sfi := sortableFileInfo(fis)
	sfi = FilterDir(sfi)
	sort.Reverse(sfi)

	recent := make([]*ShortPost, maxRecentPosts)
	all := make([]*LongPost, len(sfi))
	// First pass to get the recent posts (and others) so that
	// they can be passed to all posts.
	for i, fi := range sfi {
		all[i] = newLongPost(fi)
		if i < maxRecentPosts {
			recent[i] = all[i].Short()
		}
	}

	for i, p := range all {
		td := newTemplateData(p, i, recent, all)
		regenerateFile(td, i == 0)
	}
}

type TemplateData struct {
	Post   *LongPost
	Recent []*ShortPost
	Prev   *ShortPost
	Next   *ShortPost
}

func newTemplateData(p *LongPost, i int, r []*ShortPost, all []*LongPost) *TemplateData {
	td := &TemplateData{Post: p, Recent: r}

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

// TODO : Should pass to the template:
// Title : The first heading in the file, or the file name, or front matter?
// Description : ?
// ModTime
// Parsed : The html-parsed markdown
// Recent : A slice of n recent posts
// Next : The next (more recent) post
// Previous : The previous (older) post

func regenerateFile(td *TemplateData, idx bool) {
	var w io.Writer

	fw, err := os.Create(filepath.Join(PublicDir, td.Post.Slug))
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	defer fw.Close()
	w = fw
	if idx {
		idxw, err := os.Create(filepath.Join(PublicDir, "index.html"))
		if err != nil {
			log.Fatal("FATAL ", err)
		}
		defer idxw.Close()
		w = io.MultiWriter(fw, idxw)
	}
	err = postTpl.ExecuteTemplate(w, postTplNm, td)
	if err != nil {
		log.Fatal("FATAL ", err)
	}
}
