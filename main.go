package main

import (
	"flag"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		displayUsage()
		os.Exit(0)
	}

	isMulti := false
	withVoice := flag.Bool("v", false, "Query with speech")
	withMore := flag.Bool("m", false, "Query with more example phrases and sentence")
	openBrowser := flag.Bool("w", false, "Query and open browser to get detail")

	listWords := flag.Bool("list", false, "Query word history")
	clean := flag.Bool("clean", false, "Clean word history")
	dump := flag.Bool("dump", false, "dump query history to csv file")

	flag.Parse()

	if *listWords {
		prettyWords()
		os.Exit(0)
	}
	if *clean {
		cleanHistory()
		os.Exit(0)
	}

	if *dump {
		dumpHistory()
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		displayUsage()
		os.Exit(0)
	} else {
		words := flag.Args()
		if flag.NArg() > 1 {
			isMulti = true
		}

		query(words, *withVoice, *withMore, *openBrowser, isMulti)
	}
}
