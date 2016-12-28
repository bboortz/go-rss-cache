package main

import (
//	"github.com/op/go-logging"
	"time"
)


func logAccess(route string, method string, uri string, logTime time.Time) {
//	var log = logging.MustGetLogger("central-router")
	log.Infof( "%s\t%s\t%s\t%s", route, method, uri, time.Since(logTime) )
}

func logServiceRegistered(s Service) {
	log.Infof( "service registered: %d - %s ", s.Id, s.Name )
}
