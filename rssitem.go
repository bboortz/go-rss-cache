package main

import (
	"rsslib"
)

type RssItemCreate struct {
	Channel string `json:"Channel"`
	Title   string `json:"Title"`
}

type RssItemCreated struct {
	Item   string `json:"Service"`
	Status string `json:"Status"`
	Desc   string `json:"Desc"`
}

var rssItems rsslib.RssItems
var currentId int64 = 0

func addItem(s rsslib.RssItem) rsslib.RssItem {
	currentId += 1
	s.Id = currentId
	rssItems = append(rssItems, s)
	logItemAdded(s)
	return s
}

func findItem(title string) rsslib.RssItem {
	for _, s := range rssItems {
		if s.Title == title {
			return s
		}
	}
	// return empty Todo if not found
	return rsslib.RssItem{}
}
