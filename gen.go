package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"

	"github.com/eknkc/amber"
)

var (
	postTpl   *template.Template
	postTplNm = "post.amber"
	rssTplNm  = "rss.amber"
)

type sortableLongPost []*LongPost

func (s sortableLongPost) Len() int           { return len(s) }
func (s sortableLongPost) Less(i, j int) bool { return s[i].PubTime.Before(s[j].PubTime) }
func (s sortableLongPost) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func Filter(fi []os.FileInfo) []os.FileInfo {
	for i := 0; i < len(fi); {
		if fi[i].IsDir() || filepath.Ext(fi[i].Name()) != ".md" {
			fi[i], fi = fi[len(fi)-1], fi[:len(fi)-1]
		} else {
			i++
		}
	}
	return fi
}

func compileTemplate() error {
	ap := filepath.Join(TemplatesDir, postTplNm)
	if _, err := os.Stat(ap); os.IsNotExist(err) {
		// Amber post template does not exist, compile the native Go templates
		postTpl, err = template.ParseGlob(filepath.Join(TemplatesDir, "*.html"))
		if err != nil {
			return fmt.Errorf("error parsing templates: %s", err)
		}
		postTplNm = "post" // TODO : Validate this...
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

func clearPublicDir() error {
	// Clear the public directory, except subdirs
	fis, err := ioutil.ReadDir(PublicDir)
	if err != nil {
		return fmt.Errorf("error getting public directory files: %s", err)
	}
	for _, fi := range fis {
		if !fi.IsDir() && fi.Name() != "favicon.ico" {
			err = os.Remove(filepath.Join(PublicDir, fi.Name()))
			if err != nil {
				return fmt.Errorf("error deleting file %s: %s", fi.Name(), err)
			}
		}
	}
	return nil
}

func generateSite() error {
	// First compile the template(s)
	if err := compileTemplate(); err != nil {
		return err
	}
	// Now read the posts
	fis, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		return err
	}
	// Remove directories from the list, keep only .md files
	fis = Filter(fis)

	// Get all posts.
	all := make(sortableLongPost, len(fis))
	for i, fi := range fis {
		all[i] = newLongPost(fi)
	}
	// Then sort in reverse order (newer first)
	sort.Sort(sort.Reverse(all))
	// Slice to get only recent posts
	recent := all[:Options.RecentPostsCount]
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
