package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"compress/gzip"
)
var basedir string
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run server.go [basedir]")
		os.Exit(1)
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
	    log.Fatal(err)
	}
	
	gz, err := gzip.NewReader(f)
	if err != nil {
	    log.Fatal(err)
	}
	
	io.Copy(os.Stdout, gz)
}

