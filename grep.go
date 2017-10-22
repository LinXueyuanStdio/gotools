package main

import (
	"fmt"
	"os"
	//"io"
	"log"
	"flag"
	"bufio"
	"regexp"
)

var invertMatch = flag.Bool("v", false, "Selected lines are those not matching any of the specified patterns.")
var caseInsensitive = flag.Bool("i", false, "Perform case insensitive matching.")
var showLineNumber = flag.Bool("n", false, "Each output line is preceded by its relative line number in the file, starting at line 1.")


func main() {
    flag.Parse()
    args := flag.Args()
    
    fmt.Println("inverMatch", *invertMatch)
	fmt.Println("caseInsensitive", *caseInsensitive)
	fmt.Println("showLineNumber", *showLineNumber)
	fmt.Println("args", args)
	
    regstr := args[0]
    files  := args[1:]
    fmt.Println(regstr,"---",files)
    
    var re *regexp.Regexp
    if *caseInsensitive {
        re, _ = regexp.Compile(fmt.Sprint("(?i)",regstr))
    } else {
        re, _ = regexp.Compile(regstr)
    }
    
    if len(files) != 0 {
        var showFileName bool
        if len(files) >= 2 {
            showFileName=true
        }else{
            showFileName=false
        }
        for i:=0; i<len(files); i++ {
            f, err := os.Open(files[i])
            if err != nil {
                log.Fatal(err)
            }
            scanner := bufio.NewScanner(f)
            for j:=0;scanner.Scan();j++ {
                if re.MatchString(scanner.Text()) == !*invertMatch {
                    switch {
                    case showFileName && *showLineNumber :
                        fmt.Println(files[i], j, scanner.Text())
                    case !showFileName && *showLineNumber :
                        fmt.Println(j, scanner.Text())
                    case showFileName && !*showLineNumber :
                        fmt.Println(files[i], scanner.Text())
                    case !showFileName && !*showLineNumber :
                        fmt.Println(scanner.Text())
                    }
		        }
	        }
        }
    } else {
        scanner := bufio.NewScanner(os.Stdin)
        for i:=0;scanner.Scan();i++ {
            //i++
            if re.MatchString(scanner.Text()) == !*invertMatch {
                if *showLineNumber {
		            fmt.Println(i, scanner.Text())
		        } else {
		            fmt.Println(scanner.Text()) 
		        }
		    }
	    }
	}
}

