package main

import (
    "os"
)

func main(){
    if len(os.Args) == 1 {
        displayUsage()
        os.Exit(0)
    }

    words, withVoice, withMore, openBrowser := parseArgs(os.Args)
    query(words, withVoice, withMore, openBrowser, len(words) > 1)
}
