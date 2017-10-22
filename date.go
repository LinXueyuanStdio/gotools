package main

import (
	"time"
	"fmt"
)

func main() {
    fmt.Println(time.Now().Format(time.RFC1123))
}
