package main

import (
	"fmt"
	"os"
	"runtime"
	"unicode"

	"github.com/fatih/color"
)

var (
	version = "1.2"
	logo    = `
██╗   ██╗██████╗  █████╗  ██████╗
╚██╗ ██╔╝██╔══██╗██╔══██╗██╔═══██╗
 ╚████╔╝ ██║  ██║███████║██║   ██║
  ╚██╔╝  ██║  ██║██╔══██║██║   ██║
   ██║   ██████╔╝██║  ██║╚██████╔╝
   ╚═╝   ╚═════╝ ╚═╝  ╚═╝ ╚═════╝

ydao V%s 好好学学英语吧...

    `
)

func displayUsage() {
	logo = fmt.Sprintf(logo, version)
	color.Cyan(logo)
	color.Cyan("Usage:")
	color.Cyan("ydao <option> <word>          	Query the word(s)")
	color.Cyan("ydao -v <word(s) to query>    	Query with speech")
	color.Cyan("ydao -m <word(s) to query>    	Query with more example phrases and sentences")
	color.Cyan("ydao -w <word(s) to query>    	Query and open browser to get detail")
	color.Cyan("ydao -list			List query word histor")
	color.Cyan("ydao -clean			Clean query word histor")
	color.Cyan("ydao -dump			Dump query word histor")
}


func isChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

func isAvailablesOS() bool {
	return runtime.GOOS == "darwin" || runtime.GOOS == "linux"
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
