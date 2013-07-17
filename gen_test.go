package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"sort"
	"testing"
	"time"
)

func mustParse(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSort(t *testing.T) {
	ps := make(sortablePosts, 5)
	ps[0] = &LongPost{
		ShortPost: &ShortPost{
			Title:   "a",
			PubTime: mustParse("2012-01-07"),
		},
	}
	ps[1] = &LongPost{
		ShortPost: &ShortPost{
			Title:   "b",
			PubTime: mustParse("2012-04-22"),
		},
	}
	ps[2] = &LongPost{
		ShortPost: &ShortPost{
			Title:   "c",
			PubTime: mustParse("2012-01-01"),
		},
	}
	ps[3] = &LongPost{
		ShortPost: &ShortPost{
			Title:   "d",
			PubTime: mustParse("2011-11-30"),
		},
	}
	ps[4] = &LongPost{
		ShortPost: &ShortPost{
			Title:   "e",
			PubTime: mustParse("2012-12-01"),
		},
	}
	sort.Sort(ps)

	buf := bytes.NewBuffer(nil)
	for _, p := range ps {
		buf.WriteString(p.Title)
	}
	if buf.String() != "dcabe" {
		t.Errorf("expected 'dcabe', got %s", buf.String())
	}
}

func BenchmarkGenerateSite(b *testing.B) {
	b.StopTimer()
	log.SetOutput(ioutil.Discard)
	Options.RecentPostsCount = 5
	PublicDir = "./examples/amber/public"
	PostsDir = "./examples/amber/posts"
	TemplatesDir = "./examples/amber/templates"
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := generateSite()
		if err != nil {
			b.Fatal(err)
		}
	}
}
