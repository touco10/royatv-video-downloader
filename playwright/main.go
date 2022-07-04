package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/canhlinh/hlsdl"
	"github.com/playwright-community/playwright-go"
	"log"
	"strconv"
	"strings"
)

var liens []string
var m3u []string

func main() {
	pw, err := playwright.Run()
	checkerror(err)
	browser, err := pw.Chromium.Launch()
	checkerror(err)
	page, err := browser.NewPage()
	checkerror(err)
	_, err = page.Goto("https://roya.tv/program/1067")
	checkerror(err)
	parsedHtml, err := page.InnerHTML("#__layout > div > div.py-4.main-pages > section > div > div.row.episode-row")
	soupPage := soup.HTMLParse(parsedHtml)
	doc := soupPage.FindAll("a")
	for i, link := range doc {
		fmt.Printf("get link:%v",i)
		if i%3 == 0 {
			liens = append(liens, "https://roya.tv"+link.Attrs()["href"])
		}

		browser.Close()
	}
	for _, link2 := range liens {
		pw2, err1 := playwright.Run()
		checkerror(err1)
		browser2, err1 := pw2.Chromium.Launch()
		checkerror(err1)
		page2, err1 := browser2.NewPage()
		checkerror(err)
		_, err = page2.Goto(link2)
		checkerror(err)
		parsedHtml2, err1 := page2.InnerHTML(".video-container")
		checkerror(err1)
		soupPage2 := soup.HTMLParse(parsedHtml2)
		doc2 := soupPage2.FindAll("source")

		for _, links := range doc2 {
			video := (strings.Replace(links.Attrs()["src"], "playlist", "hlssubplaylist-video_1800000-audio_128968", 1))
			m3u = append(m3u, video)

		}

		browser.Close()

	}
	download()
}

//
func download() {

	for i, link := range m3u {
		fmt.Printf("downloading file %v",i)
		video := "download/saison_4/episode_" + strconv.Itoa(len(m3u)-i)
		hlsDL := hlsdl.New(link, nil, video, 64, true)
		filepath, err := hlsDL.Download()
		if err != nil {
			panic(err)
		}

		fmt.Println(filepath)
	}
}
func checkerror(err error) {
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
}
