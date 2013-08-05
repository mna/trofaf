package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/amber"
)

var (
	postTpl   *template.Template // The one and only compiled post template
	postTplNm = "post.amber"     // The amber post template file name (native Go are compiled using ParseGlob)

	// Special files in the public directory, that must not be deleted
	specFiles = map[string]struct{}{
		"favicon.ico":                              struct{}{},
		"robots.txt":                               struct{}{},
		"humans.txt":                               struct{}{},
		"crossdomain.xml":                          struct{}{},
		"apple-touch-icon.png":                     struct{}{},
		"apple-touch-icon-114x114-precomposed.png": struct{}{},
		"apple-touch-icon-144x144-precomposed.png": struct{}{},
		"apple-touch-icon-57x57-precomposed.png":   struct{}{},
		"apple-touch-icon-72x72-precomposed.png":   struct{}{},
		"apple-touch-icon-precomposed.png":         struct{}{},
	}

	funcs = template.FuncMap{
		"fmttime": func(t time.Time, f string) string {
			return t.Format(f)
		},
	}
)

func init() {
	// Add the custom functions to Amber in the init(), since this is global
	// (package) state in my Amber fork.
	amber.AddFuncs(funcs)
}

// This type is a slice of *LongPost that implements the sort.Interface, to sort in PubTime order.
type sortablePosts []*LongPost

func (s sortablePosts) Len() int           { return len(s) }
func (s sortablePosts) Less(i, j int) bool { return s[i].PubTime.Before(s[j].PubTime) }
func (s sortablePosts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Filter cleans the slice of FileInfo to leave only `.md` files (markdown)
func filter(fi []os.FileInfo) []os.FileInfo {
	for i := 0; i < len(fi); {
		if fi[i].IsDir() || filepath.Ext(fi[i].Name()) != ".md" {
			fi[i], fi = fi[len(fi)-1], fi[:len(fi)-1]
		} else {
			i++
		}
	}
	return fi
}

// Compile the Post template.
func compileTemplates() error {
	ap := filepath.Join(TemplatesDir, postTplNm)
	if _, err := os.Stat(ap); os.IsNotExist(err) {
		// Amber post template does not exist, compile the native Go templates
		postTpl, err = template.New("templates").Funcs(funcs).ParseGlob(filepath.Join(TemplatesDir, "*.html"))
		if err != nil {
			return fmt.Errorf("error parsing templates: %s", err)
		}
		postTplNm = "post"
	} else {
		c := amber.New()
		if err := c.ParseFile(ap); err != nil {
			return fmt.Errorf("error parsing templates: %s", err)
		}
		if postTpl, err = c.Compile(); err != nil {
			return fmt.Errorf("error compiling templates: %s", err)
		}
	}
	return nil
}

// Clear the public directory, ignoring special files, subdirectories, and hidden (dot) files.
func clearPublicDir() error {
	// Clear the public directory, except subdirs and special files (favicon.ico & co.)
	fis, err := ioutil.ReadDir(PublicDir)
	if err != nil {
		return fmt.Errorf("error getting public directory files: %s", err)
	}
	for _, fi := range fis {
		if !fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			// Check for special files
			if _, ok := specFiles[fi.Name()]; !ok {
				err = os.Remove(filepath.Join(PublicDir, fi.Name()))
				if err != nil {
					return fmt.Errorf("error deleting file %s: %s", fi.Name(), err)
				}
			}
		}
	}
	return nil
}

func getPosts(fis []os.FileInfo) (all, recent []*LongPost) {
	all = make([]*LongPost, 0, len(fis))
	for _, fi := range fis {
		lp, err := newLongPost(fi)
		if err == nil {
			all = append(all, lp)
		} else {
			log.Printf("post ignored: %s; error: %s\n", fi.Name(), err)
		}
	}
	// Then sort in reverse order (newer first)
	sort.Sort(sort.Reverse(sortablePosts(all)))
	cnt := Options.RecentPostsCount
	if l := len(all); l < cnt {
		cnt = l
	}
	// Slice to get only recent posts
	recent = all[:cnt]
	return
}

// Generate the whole site.
func generateSite() error {
	// First compile the template(s)
	if err := compileTemplates(); err != nil {
		return err
	}
	// Now read the posts
	fis, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		return err
	}
	// Remove directories from the list, keep only .md files
	fis = filter(fis)
	// Get all posts.
	all, recent := getPosts(fis)
	// Delete current public directory files
	if err := clearPublicDir(); err != nil {
		return err
	}
	// Generate the static files
	for i, p := range all {
		td := newTemplateData(p, i, recent, all)
		if err := generateFile(td, i == 0); err != nil {
			return err
		}
	}
	// Generate the RSS feed
	td := newTemplateData(nil, 0, recent, nil)
	return generateRss(td)
}

// Creates the rss feed from the recent posts.
func generateRss(td *TemplateData) error {
	r := NewRss(td.SiteName, td.TagLine, Options.BaseURL)
	base, err := url.Parse(Options.BaseURL)
	if err != nil {
		return fmt.Errorf("error parsing base URL: %s", err)
	}
	for _, p := range td.Recent {
		u, err := base.Parse(p.Slug)
		if err != nil {
			return fmt.Errorf("error parsing post URL: %s", err)
		}
		r.Channels[0].AppendItem(NewRssItem(p.Title, u.String(), p.Description, p.Author, "", p.PubTime))
	}
	return r.WriteToFile(filepath.Join(PublicDir, "rss"))
}

// Generate the static HTML file for the post identified by the index.
func generateFile(td *TemplateData, idx bool) error {
	var w io.Writer

	fw, err := os.Create(filepath.Join(PublicDir, td.Post.Slug))
	if err != nil {
		return fmt.Errorf("error creating static file %s: %s", td.Post.Slug, err)
	}
	defer fw.Close()

	// If this is the newest file, also save as index.html
	w = fw
	if idx {
		idxw, err := os.Create(filepath.Join(PublicDir, "index.html"))
		if err != nil {
			return fmt.Errorf("error creating static file index.html: %s", err)
		}
		defer idxw.Close()
		w = io.MultiWriter(fw, idxw)
	}
	return postTpl.ExecuteTemplate(w, postTplNm, td)
}
