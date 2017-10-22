package main

import (
	"io"
	"log"
	"os"
	"fmt"
)

func main() {
    count := len(os.Args)
    for i := 1 ; i < count ;i++{  
        file, err := os.Open(os.Args[i]) // For read access.
        if err != nil {
            log.Fatal(err)
            continue
        }
        if _, err := io.Copy(os.Stdout, file); err != nil {
		    log.Fatal(err)
	    }
        file.Close()
    }
}
