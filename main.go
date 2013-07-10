package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/eknkc/amber"
	"github.com/howeyc/fsnotify"
	"github.com/jessevdk/go-flags"
)

var (
	postTpl   *template.Template
	postTplNm = "post.amber"
)

func main() {
	_, err := flags.Parse(&Options)
	if err == nil { // err prints the usage automatically
		// Compile the template(s)
		compileTemplate()
		// Generate the site
		regeneratePosts()

		// Start the watcher
		w, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal("FATAL ", err)
		}
		defer w.Close()
		go watch(w)
		if err = w.Watch(PostsDir); err != nil {
			log.Fatal("FATAL ", err)
		}

		// Start the web server
		run()
	}
}

func compileTemplate() {
	ap := filepath.Join(TemplatesDir, postTplNm)
	if _, err := os.Stat(ap); os.IsNotExist(err) {
		// Amber post template does not exist, compile the native Go templates
		postTpl, err = template.ParseGlob(filepath.Join(TemplatesDir, "*.html"))
		if err != nil {
			log.Fatal("FATAL ", err)
		}
		postTplNm = "post" // TODO : Validate this...
	} else {
		c := amber.New()
		if err := c.ParseFile(ap); err != nil {
			log.Fatal("FATAL ", err)
		}
		if postTpl, err = c.Compile(); err != nil {
			log.Fatal("FATAL ", err)
		}
	}
}
