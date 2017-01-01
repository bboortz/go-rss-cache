package main

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("rss-cache")

func logServiceRegistered(s Service) {
	log.Infof("%s\tservice registered: %d - %s ", log.Module, s.Id, s.Name)
}
