package main

import (
	"os"
	"runtime"
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
	



// PrintStack prints to standard error the stack trace returned by runtime.Stack.
func printStack() {
	os.Stderr.Write(Stack())
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
