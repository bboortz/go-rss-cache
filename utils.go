package main

import (
	"os"
	"reflect"
	"runtime"
)


func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
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

