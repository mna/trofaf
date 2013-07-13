package main

// rss_api.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

/*
Package to parse and create RSS-Feeds
*/

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"time"
)

type Rss struct {
	XMLName xml.Name  `xml:"rss"`
	Version string    `xml:"version,attr"`
	Channel []Channel `xml:"channel"`
}

type Channel struct {
	Title         string  `xml:"title"`
	Description   string  `xml:"description"`
	Link          string  `xml:"link"`
	LastBuildDate string  `xml:"lastBuildDate"`
	Generator     string  `xml:"generator"`
	Image         []Image `xml:"image"`
	Item          []Item  `xml:"item"`
}

type Image struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Item struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	Author      string  `xml:"author"`
	Category    string  `xml:"category"`
	PupDate     string  `xml:"pubDate"`
	Image       []Image `xml:"image"`
}

func ParseFromFile(filename string) (*Rss, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseFromReader(file)
}

func ParseFromUrl(url string) (*Rss, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseFromReader(resp.Body)
}

func ParseFromReader(reader io.Reader) (*Rss, error) {
	var rss Rss
	dec := xml.NewDecoder(reader)
	err := dec.Decode(&rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

func New(title string, description string, link string) *Rss {
	rss := &Rss{Version: "2.0",
		Channel: []Channel{Channel{
			Title:       title,
			Description: description,
			Link:        link,
			Generator:   "gbt",
			Image:       make([]Image, 0),
			Item:        make([]Item, 0)}}}

	return rss
}

// Add a new Item to the feed
func (rss *Rss) AddItem(title string, link string, description string, author string, category string) {
	item := Item{
		Title:       title,
		Link:        link,
		Description: description,
		Author:      author,
		Category:    category,
		PupDate:     time.Now().Format(time.RFC822Z),
		Image:       make([]Image, 0)}

	//prepend s = append(s, T{}); copy(s[1:], s); s[0] = prefix
	rss.Channel[0].Item = append(rss.Channel[0].Item, Item{})
	copy(rss.Channel[0].Item[1:], rss.Channel[0].Item)
	rss.Channel[0].Item[0] = item
}

// Writes the data in RSS 2.0 format to a given file
func (rss *Rss) WriteToFile(path string) error {
	rss.Channel[0].LastBuildDate = time.Now().Format(time.RFC822Z)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := xml.NewEncoder(file)

	return enc.Encode(rss)
}
