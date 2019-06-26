// 爬取豆瓣小组
package main

import (
	"log"

	//"strings"

	"github.com/PuerkitoBio/goquery"

	"go-crawler/douban-group/model"
	"go-crawler/douban-group/parse"
)

// 抓取网站：豆瓣🔥小组
var (
	BaseURL = "https://www.douban.com/group/639264/discussion"
)

// Start 开始爬取
func Start() {
	var topics []parse.DoubanGroupDbhyz

	lastUpdate, err := model.GetLastTime()
	if err != nil {
		log.Println(err)
		return
	}

	var pages []parse.Page
	if lastUpdate != "" {
		pages = parse.GetPages(BaseURL, 2)
	} else {
		pages = parse.GetPages(BaseURL, 0)
	}

	for _, page := range pages {
		doc, err := goquery.NewDocument(page.URL)
		if err != nil {
			log.Println(err)
		}

		topics = append(topics, parse.Topics(doc)...)
	}

	model.Save(topics)
}

func main() {
	Start()

	defer model.DB.Close()
}
