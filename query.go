package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "time"
    "os"
    "os/exec"
    "net/http"
    "strings"

	"github.com/briandowns/spinner"
	"github.com/PuerkitoBio/goquery"
    "github.com/fatih/color"
)

const (
    voiceURL = "https://dict.youdao.com/dictvoice?audio=%s&type=2"
    enURL = "http://dict.youdao.com/w/eng/%s"
    znURL = "http://dict.youdao.com/w/%s"
    multiURL = "http://dict.youdao.com/example/blng/eng/%s"
)

func query(words []string, withVoice, withMore, isMulti bool) {
    var url string
    var doc *goquery.Document
    var voiceBody io.ReadCloser

    queryString := strings.Join(words, " ")
	voiceString := strings.Join(words, "+")

    isChinese := isChinese(queryString)
    if url = enURL; isChinese {
        url = znURL
    }

    s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Prefix = "Querying... "
	s.Color("green")
	s.Start()

    var err error
    doc, err = goquery.NewDocument(fmt.Sprintf(url, queryString))

    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    if withVoice && isAvailablesOS() {
        if resp, err := http.Get(fmt.Sprintf(voiceURL, voiceString)); err == nil {
            voiceBody = resp.Body
        }
    }

    s.Stop()

    if isChinese {
        // Find the result
        fmt.Println()
        doc.Find(".trans-container > ul > p > span.contentTitle").Each(func (i int, s *goquery.Selection) {
            color.Green("   %s", s.Find(".search-js").Text())
        })
    } else {
        // Check for typos
        if hint := getHint(doc); hint != nil {
            color.Blue("\r\n   word '%s' not found, do you mean?", queryString)
            fmt.Println()
            for _, guess := range hint {
                color.Green("   %s", guess[0])
                color.Magenta("  %s", guess[1])
            }
            fmt.Println()
            return
        }

        // Find the pronounce
        if !isMulti {
            color.Yellow("音标：")
            color.Green("\r\n   %s", getPronounce(doc))
        }

        // Find the result
	    result := doc.Find("div#phrsListTab > div.trans-container > ul").Text()
	    color.Green(result)
     }

    // get Phrase
    phrase := getPhrase(doc, withMore)
    if len(phrase) > 0 {
        fmt.Println()
        color.Yellow("短语：")
        for i, _ := range phrase {
            color.Green("   %s", phrase[i])
        }
    }

    // Show examples
    sentences := getSentences(words, doc, isChinese, withMore)
    if len(sentences) > 0 {
        fmt.Println()
        color.Yellow("句子：")
        for i, sentence := range sentences {
            color.Green(" %2d.%s", i+1, sentence[0])
            color.Magenta("     %s", sentence[1])
        }
        fmt.Println()
    }

    // Play voice
    if withVoice && isAvailablesOS() {
        playVoice(voiceBody)
    }
}

func getHint(doc *goquery.Document) [][]string {
    typos := doc.Find(".typo-rel")
    if typos.Length() == 0 {
        return nil
    }

    results := [][]string{}
    typos.Each(func(_ int, s *goquery.Selection) {
        word := strings.TrimSpace(s.Find("a").Text())
        s.Children().Remove()
        mean := strings.TrimSpace(s.Text())
        results = append(results, []string{word, mean})
    })
    return results
}

// 获取音标，第一个元素英式，第二为美式
func getPronounce(doc *goquery.Document) string {
    var pronouce, pronouceType string

    doc.Find("div.baav > span.pronounce").Each(func (i int, s *goquery.Selection) {
        phonetic := s.Find("span.phonetic").Text()
        if pronouceType = "英: "; i == 1 {
            pronouceType = "美：";
        }
        pronouce += (pronouceType + phonetic + "    ")
    })
    return pronouce
}

func getPhrase(doc *goquery.Document, withMore bool) []string {
    result := []string{}
    doc.Find("div#webPhrase > p").Each(func (i int, s *goquery.Selection) {
        phraseDom := s.Find("a.search-js")
        phrase := phraseDom.Text()
        phraseDom.Remove()
        txt := prettyPhrase(strings.TrimSpace(s.Text()))
        cnt := phrase + "   " + txt
        result = append(result, cnt)
    })
    if !withMore {
        return result[:3]
    }
    return result
}

func getSentences(words []string, doc *goquery.Document, isChinese, withMore bool) [][]string {
    result := [][]string{}
    if withMore {
        url := fmt.Sprintf(multiURL, strings.Join(words, "-"))
        var err error
        doc, err = goquery.NewDocument(url)
        if err != nil {
            return result
        }
    }

    doc.Find("#bilingual ul li").Each(func(_ int, s *goquery.Selection) {
        r := []string{}
        s.Children().Each(func (ii int, ss *goquery.Selection) {
            // ignore source
            if ii == 2 {
                return
            }

            var sentence string
            ss.Children().Each(func (iii int, sss *goquery.Selection) {
                if text := strings.TrimSpace(sss.Text()); text != "" {
                    addSpace := (ii == 1 && isChinese) || (ii ==0 && !isChinese) && iii != 0 && text != "."
                    if addSpace {
                        text = " " + text
                    }
                    sentence += text
                }
            })
            r = append(r, sentence)
        })
        if len(r) == 2 {
            result = append(result, r)
        }
    })
    return result
}

// TODO: open browser
// TODO: save history to sqlite

func playVoice(body io.ReadCloser) {
    tmpfile, err := ioutil.TempFile("", "ydict")
    if err != nil {
        log.Fatal(err)
    }

    // clean up
    defer os.Remove(tmpfile.Name())

    data, err := ioutil.ReadAll(body)
    if err != nil {
        log.Fatal(err)
    }

    if _, err := tmpfile.Write(data); err != nil {
        log.Fatal(err)
    }

    cmd := exec.Command("mpg123", tmpfile.Name())

    if err := cmd.Start(); err != nil {
        fmt.Println(err)
    }

    if err := cmd.Wait(); err != nil {
        fmt.Println(err)
    }
}

func prettyPhrase(text string) string {
    ss := strings.Split(text, ";")
    for i, m := range ss {
        ss[i] = strings.TrimSpace(m)
    }

   return strings.Join(ss, ";")
}
