package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/ghost/handlers"
)

var (
	// Favicon path and cache duration
	faviconPath  = filepath.Join(PublicDir, "favicon.ico")
	faviconCache = 2 * 24 * time.Hour
)

// Start serving the blog.
func run() {
	h := handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				handlers.GZIPHandler(
					http.FileServer(http.Dir(PublicDir)),
					nil),
				handlers.NewLogOptions(nil, handlers.Lshort)),
			nil),
		faviconPath,
		faviconCache)

	// Assign the combined handler to the server.
	http.Handle("/", h)

	// Start it up.
	log.Printf("trofaf server listening on port %d", Options.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Options.Port), nil); err != nil {
		log.Fatal("FATAL ", err)
	}
}
