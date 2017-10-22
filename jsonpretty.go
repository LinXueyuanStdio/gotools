package main

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
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
	
	var data interface{}
	
	decoder := json.NewDecoder(f)
    err = decoder.Decode(&data)
    log.Printf("decoded data: %v", data)
    
    saver, err := os.Create("save.json")
    if err != nil {
	    log.Fatal(err)
	}
    encoder := json.NewEncoder(os.Stdout) //f不行,saver不行,os.Stdout也不行
    /*
    # command-line-arguments
    ./jsonpretty.go:29: encoder.SetIndent undefined (type *json.Encoder has no field or method SetIndent)

    */
	encoder.SetIndent("", "  ")
	
	err = encoder.Encode(data)
	
	if err != nil {
	    log.Fatal(err)
	}
	
}

