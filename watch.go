package main

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
)

const (
	// Sometimes many events can be triggered in succession for the same file
	// (i.e. Create followed by Modify, etc.). No need to rush to generate
	// the HTML, just wait for it to calm down before processing.
	watchEventDelay = 10 * time.Second
)

// Create and start a watcher, watching both the posts and the templates directories.
func startWatcher() *fsnotify.Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	go watch(w)
	// Watch the posts directory
	if err = w.Watch(PostsDir); err != nil {
		w.Close()
		log.Fatal("FATAL ", err)
	}
	// Watch the templates directory
	if err = w.Watch(TemplatesDir); err != nil {
		w.Close()
		log.Fatal("FATAL ", err)
	}
	return w
}

// Receive watcher events for the directories. All events require re-generating
// the whole site (because the template may display the n most recent posts, the
// next and previous post, etc.). It could be fine-tuned based on what data we give
// to the templates, but for now, lazy approach.
func watch(w *fsnotify.Watcher) {
	var delay <-chan time.Time
	for {
		select {
		case ev := <-w.Event:
			// Regenerate the files after the delay, reset the delay if an event is triggered
			// in the meantime
			ext := filepath.Ext(ev.Name)
			// Care only about changes to markdown files in the Posts directory, or to
			// Amber or Native Go template files in the Templates directory.
			if strings.HasPrefix(ev.Name, PostsDir) && ext == ".md" {
				delay = time.After(watchEventDelay)
			} else if strings.HasPrefix(ev.Name, TemplatesDir) && (ext == ".amber" || ext == ".html") {
				delay = time.After(watchEventDelay)
			}

		case err := <-w.Error:
			log.Println("WATCH ERROR ", err)

		case <-delay:
			if err := generateSite(); err != nil {
				log.Println("ERROR generating site: ", err)
			} else {
				log.Println("site generated")
			}
		}
	}
}
