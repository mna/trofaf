package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Port             int    `short:"p" long:"port" description:"the port to use for the web server" default:"9000"`
	NoGen            bool   `short:"G" long:"no-generation" description:"when set, the site is not automatically generated"`
	SiteName         string `short:"n" long:"site-name" description:"the name of the site" default:"Site Name"`
	TagLine          string `short:"t" long:"tag-line" description:"the site's tag line"`
	RecentPostsCount int    `short:"r" long:"recent-posts" description:"the number of recent posts to send to the templates" default:"5"`
	BaseURL          string `short:"b" long:"base-url" description:"the base URL of the web site" default:"http://localhost"`
}

var (
	Options      options
	PublicDir    string
	PostsDir     string
	TemplatesDir string
)

func main() {
	// Initialize directories
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	PublicDir = filepath.Join(pwd, "public")
	PostsDir = filepath.Join(pwd, "posts")
	TemplatesDir = filepath.Join(pwd, "templates")

	// Parse the flags
	_, err = flags.Parse(&Options)
	if err == nil { // err != nil prints the usage automatically
		if !Options.NoGen {
			// Generate the site
			if err := generateSite(); err != nil {
				log.Fatal("FATAL ", err)
			}
			// Start the watcher
			defer startWatcher().Close()
		}
		// Start the web server
		run()
	}
}
