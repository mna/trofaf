package main

// Adapted from https://github.com/krautchan/gbt
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
	"os"
	"time"
)

// The root Rss structure
type Rss struct {
	XMLName  xml.Name   `xml:"rss"`
	Version  string     `xml:"version,attr"`
	Channels []*Channel `xml:"channel"`
}

// The Rss channel structure
type Channel struct {
	Title         string   `xml:"title"`
	Description   string   `xml:"description"`
	Link          string   `xml:"link"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Generator     string   `xml:"generator"`
	Image         []*Image `xml:"image"`
	Item          []*Item  `xml:"item"`
}

// The rss image structure
type Image struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

// The Rss item structure
type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Author      string   `xml:"author"`
	Category    string   `xml:"category"`
	PubDate     string   `xml:"pubDate"`
	Image       []*Image `xml:"image"`
}

// Create a new RSS feed
func NewRss(title string, description string, link string) *Rss {
	rss := &Rss{Version: "2.0",
		Channels: []*Channel{
			&Channel{
				Title:       title,
				Description: description,
				Link:        link,
				Generator:   "trofaf (https://github.com/PuerkitoBio/trofaf)",
				Image:       make([]*Image, 0),
				Item:        make([]*Item, 0),
			},
		},
	}

	return rss
}

// Create a new, orphan Rss Item.
func NewRssItem(title, link, description, author, category string, pubTime time.Time) *Item {
	return &Item{
		Title:       title,
		Link:        link,
		Description: description,
		Author:      author,
		Category:    category,
		PubDate:     pubTime.Format(time.RFC822),
		Image:       make([]*Image, 0),
	}
}

// Add an Item to the feed, under this Channel
func (ch *Channel) AppendItem(i *Item) {
	ch.Item = append(ch.Item, i)
}

// Writes the data in RSS 2.0 format to a given file
func (rss *Rss) WriteToFile(path string) error {
	rss.Channels[0].LastBuildDate = time.Now().Format(time.RFC822)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}
	enc := xml.NewEncoder(file)
	return enc.Encode(rss)
}
