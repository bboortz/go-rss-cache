package main

import (
	"github.com/bboortz/go-rsslib"
	//"go-rsslib"
	//	"github.com/davecgh/go-spew/spew"
	"strconv"
	"time"
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

func getItem(uuid string) rsslib.RssItem {
	_, result := findItem(uuid)
	return result
}

func addItem(s rsslib.RssItem) rsslib.RssItem {
	if s.PublishDate == "" {
		s.PublishDate = time.Now().UTC().String()
	}

	s.Id = currentId
	rssItems = append(rssItems, s)
	currentId += 1
	logItemAdded(s)
	return s
}

func updateItem(i uint64, s rsslib.RssItem) rsslib.RssItem {
	if s.PublishDate == "" {
		s.PublishDate = rssItems[i].PublishDate
	}
	if s.UpdateDate == "" {
		s.UpdateDate = time.Now().UTC().String()
	}

	s.Id = i
	rssItems[i] = s
	logItemUpdated(s)
	return s
}

func addOrUpdateItem(s rsslib.RssItem) ItemCUDResult {
	var result ItemCUDResult
	var resultItem rsslib.RssItem

	searchKey, searchItem := findItem(s.Uuid)
	if searchItem.Uuid == "" {
		resultItem = addItem(s)
		result = ItemCUDResult{Item: resultItem.Uuid, Status: "created", Desc: strconv.FormatUint(resultItem.Id, 10)}
	} else {
		s.Id = searchItem.Id
		if s.UpdateDate == "" {
			s.UpdateDate = searchItem.UpdateDate
		}
		s.Diff(searchItem)
		if !s.Compare(searchItem) {
			resultItem = updateItem(uint64(searchKey), s)
			result = ItemCUDResult{Item: resultItem.Uuid, Status: "updated", Desc: strconv.FormatUint(resultItem.Id, 10)}
		} else {
			result = ItemCUDResult{Item: searchItem.Uuid, Status: "notmodified", Desc: strconv.FormatUint(resultItem.Id, 10)}
		}
	}

	return result
}

func findItem(uuid string) (int, rsslib.RssItem) {
	for k, s := range rssItems {
		if s.Uuid == uuid {
			return k, s
		}
	}
	// return empty item if not found
	return 0, rsslib.RssItem{}
}
