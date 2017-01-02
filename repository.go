package main

import (
	"rsslib"
)

type ItemCount struct {
	Count uint64 `json:"Count"`
}

type ItemCUDResult struct {
	Item   string `json:"Item"`
	Status string `json:"Status"`
	Desc   string `json:"Desc"`
}

var rssItems rsslib.RssItems
var currentId uint64 = 0

func addItem(s rsslib.RssItem) rsslib.RssItem {
	currentId += 1
	s.Id = currentId
	rssItems = append(rssItems, s)
	logItemAdded(s)
	return s
}

func findItem(uuid string) rsslib.RssItem {
	for _, s := range rssItems {
		if s.Uuid == uuid {
			return s
		}
	}
	// return empty item if not found
	return rsslib.RssItem{}
}
