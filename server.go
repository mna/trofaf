package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/ghost/handlers"
)

// Start serving the blog.
func run() {
	var (
		faviconPath  = filepath.Join(PublicDir, "favicon.ico")
		faviconCache = 2 * 24 * time.Hour
	)

	h := handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				handlers.GZIPHandler(
					http.FileServer(http.Dir(PublicDir)),
					nil),
				handlers.NewLogOptions(nil, handlers.Ldefault)),
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
