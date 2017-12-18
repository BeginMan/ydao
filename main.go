package main

import (
    "os"
)

var (
    proxy string
)

func main(){
    if len(os.Args) == 1 {
        displayUsage()
        os.Exit(0)
    }

    words, withVoice, withMore := parseArgs(os.Args)
    query(words, withVoice, withMore, len(words) > 1)
}
