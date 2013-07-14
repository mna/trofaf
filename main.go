package main

import (
	"github.com/jessevdk/go-flags"
	"log"
)

func main() {
	_, err := flags.Parse(&Options)
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
