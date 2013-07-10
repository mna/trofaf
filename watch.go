package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
)

const (
	// Sometimes many events can be triggered in succession for the same file
	// (i.e. Create followed by Modify, etc.). No need to rush to generate
	// the HTML, just wait for it to calm down before processing.
	watchEventDelay = 30 * time.Second
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
	fis, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	sfi := sortableFileInfo(fis)
	sfi = FilterDir(sfi)
	sort.Sort(sfi)
	for _, fi := range sfi {
		regenerateFile(fi)
	}
}

// TODO : Should pass to the template:
// Title : The first heading in the file, or the file name, or front matter?
// Description : ?
// ModTime
// Parsed : The html-parsed markdown
// Recent : A slice of n recent posts
// Next : The next (more recent) post
// Previous : The previous (older) post

func regenerateFile(fi os.FileInfo) {
	f, err := os.Open(filepath.Join(PostsDir, fi.Name()))
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	defer f.Close()
	// TODO : Blackfriday...

	nm := fi.Name()
	nm = strings.Replace(nm, filepath.Ext(nm), "", 1)
	fw, err := os.Create(filepath.Join(PublicDir, nm))
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	defer fw.Close()
	postTpl.ExecuteTemplate(fw, "post", nil)
}
