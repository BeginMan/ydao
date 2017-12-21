package main

import (
    "fmt"
    "strings"
    "unicode"
    "runtime"

	"github.com/fatih/color"
)

var (
    version = "1.1"
	logo = "ydao V%s 好好学学英语吧..."
)

func displayUsage() {
    logo = fmt.Sprintf(logo, version)
    color.Cyan(logo)
    color.Cyan("Usage:")
    color.Cyan("ydict <word(s) to query>        Query the word(s)")
    color.Cyan("ydict <word(s) to query> -v     Query with speech")
    color.Cyan("ydict <word(s) to query> -m     Query with more example phrases and sentences")
    color.Cyan("ydict <word(s) to query> -w     Query and open browser to get detail")
}


func parseArgs(args []string) ([]string, bool, bool, bool) {
    // match argument: -v,-m, -w
    var withVoice, withMore, openBrowser bool
    parameterStartIndex := findParaStartIndex(args)
	paramArray := args[parameterStartIndex:]

    if elementInStringArray(paramArray, "-m") {
		withMore = true
	}

	if elementInStringArray(paramArray, "-v") {
		withVoice = true
	}

    if elementInStringArray(paramArray, "-w") {
        openBrowser = true
    }
    return args[1:parameterStartIndex], withVoice, withMore, openBrowser
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
