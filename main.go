package main

import (
	"github.com/jessevdk/go-flags"
)

func main() {
	_, err := flags.Parse(&Options)
	if err == nil { // err != nil prints the usage automatically
		if !Options.NoGen {
			// Generate the site
			generateSite()
			// Start the watcher
			defer startWatcher().Close()
		}
		// Start the web server
		run()
	}
}
