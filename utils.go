package main

import (
    "strings"
    "unicode"
    "runtime"

	"github.com/fatih/color"
)

var (
	Version = "0.1"
	logo    = `
██╗   ██╗██████╗ ██╗ ██████╗████████╗
╚██╗ ██╔╝██╔══██╗██║██╔════╝╚══██╔══╝
 ╚████╔╝ ██║  ██║██║██║        ██║
  ╚██╔╝  ██║  ██║██║██║        ██║
   ██║   ██████╔╝██║╚██████╗   ██║
   ╚═╝   ╚═════╝ ╚═╝ ╚═════╝   ╚═╝


  ydict v%s
  好好学学英语吧...
    `
)

func displayUsage() {
	color.Cyan(logo, Version)
    color.Cyan("Usage:")
    color.Cyan("ydict <word(s) to query>        Query the word(s)")
    color.Cyan("ydict <word(s) to query> -v     Query with speech")
    color.Cyan("ydict <word(s) to query> -m     Query with more example phrases and sentences")
}


func parseArgs(args []string) ([]string, bool, bool) {
    // match argument: -v or -m
    var withVoice, withMore bool
    parameterStartIndex := findParaStartIndex(args)
	paramArray := args[parameterStartIndex:]

    if elementInStringArray(paramArray, "-m") {
		withMore = true
	}

	if elementInStringArray(paramArray, "-v") {
		withVoice = true
	}
    return args[1:parameterStartIndex], withVoice, withMore
}

func findParaStartIndex(args []string) int {
    for index, word := range(args) {
        if strings.HasPrefix(word, "-") && len(word) == 2 {
            return index
        }
    }
    return len(args)
}

func elementInStringArray(stringArray []string, element string) bool {
    for _, word := range(stringArray) {
        if word == element {
            return true
        }
    }
    return false
}

func isChinese (str string) bool {
    for _, r := range(str) {
        if unicode.Is(unicode.Scripts["Han"], r) {
            return true
        }
    }
    return false
}

func isAvailablesOS() bool {
    return runtime.GOOS == "darwin" || runtime.GOOS == "linux"
}
