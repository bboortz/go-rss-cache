package main

import (
	"github.com/op/go-logging"
	"rsslib"
)

var log = logging.MustGetLogger("rss-cache")

func logItemAdded(s rsslib.RssItem) {
	log.Infof("%s\titem added: %d - %s ", log.Module, s.Id, s.Title)
}
