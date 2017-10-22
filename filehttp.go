package main

import (
	"fmt"
	"os"
	"net/http"
	"io"
	"path"
)
var basedir string
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run server.go [basedir]")
		os.Exit(1)
	}

	basedir := os.Args[1]
	fmt.Printf("serving content: %s\n", basedir)

	http.HandleFunc("/", serveContent)

	http.ListenAndServe(":7777", nil)
}

func serveContent(res http.ResponseWriter, req *http.Request) {
  // TODO: req.URL.Path is the file you want to output as response
  // Hint: Use io.Copy
    fmt.Printf("serving content: %s\n", basedir) //空
    p := path.Join(basedir, "index.html")　　
    fmt.Printf("serving content: %s\n", p) 
    //index.html,应为public/index.html
    
    f, err:= os.Open(p)
    if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}
    _, err = io.Copy(res, f)
    if err != nil {
        fmt.Fprint(res, fmt.Sprint("error:",err))
    }
}
