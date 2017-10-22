package main

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"fmt"
	"strconv"
	"log"
)

func main() {
	length, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, length)
	_, err = rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	// The slice should now contain random bytes instead of only zeroes.
	
    encoded := base64.StdEncoding.EncodeToString(b)
	fmt.Println(encoded)
}
