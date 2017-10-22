package main

import (
	"io"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"os"
	"fmt"
	"strings"
)

func main() {
    var i string
    i="b"
    fmt.Println(fmt.Scan(&i)) //接受echo传的数据
    fmt.Println(i)
    var h hash.Hash
    switch os.Args[1] {
    case "md5":
        h = md5.New()
    case "sha256":
        h = sha256.New()
    case "sha512":
        h = sha512.New()
    }
    fmt.Println(*strings.NewReader(i))
    io.Copy(h, strings.NewReader(i))
    //io.WriteString(h, i)
    //io.Copy(h, os.Stdin)
    fmt.Println(*os.Stdin)
    
    var digest []byte
    digest = h.Sum(digest)
    fmt.Println(hex.EncodeToString(digest))
}
