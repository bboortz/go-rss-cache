package main

import (
	"github.com/bboortz/go-rsslib"
	"github.com/op/go-logging"
	//"go-rsslib"
)

var log = logging.MustGetLogger("rss-cache")

func logItemAdded(s rsslib.RssItem) {
	logItemAction(s, "added")
}

func logItemUpdated(s rsslib.RssItem) {
	logItemAction(s, "updated")
}

func logItemAction(s rsslib.RssItem, action string) {
	log.Infof("%s\titem %s: %d - %s - %s - %s ", log.Module, action, s.Id, s.Uuid, s.Channel, s.Title)
}
