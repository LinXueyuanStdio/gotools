package main

import (
	"os"
	"fmt"
	"log"
)

func main() {
    if os.Args[1] == "" {
        return
    }
    filename := os.Args[1]
    du(filename)
    fmt.Println("disk usage:", diskUsage)
}
var diskUsage int64
func du(filename string) {
    f, err := os.Open(filename) // For read access.
    if err != nil {
        log.Fatal(err)
        return
    }
    fileInfo, err := os.Stat(filename)
    if err != nil {
        log.Fatal(err)
        return
    }
    if fileInfo.IsDir() {
        fnset, err := f.Readdirnames(0)
        if err != nil {
            log.Fatal(err)
            return
        }
        for i:=0; i<len(fnset); i++ { //非空文件夹
            du(fmt.Sprintf("%s/%s", filename, fnset[i]))
        }
    } else { // 文件
        diskUsage = diskUsage + fileInfo.Size()
        fmt.Println(filename, fileInfo.Size())
    }
}






