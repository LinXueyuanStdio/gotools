package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"net/http"
)
var basedir string
func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: go run server.go [basedir]")
		os.Exit(1)
	}

	url := os.Args[1]
	filename := os.Args[2]
	resp, err := http.Get(url)
    if err != nil {
	    log.Fatal(err)
    }
    defer resp.Body.Close()
	f, err := os.Create(filename)
	if err != nil {
	    log.Fatal(err)
	}
	
	io.Copy(f, resp.Body)
	fmt.Println("Success!")
}

