// 爬取豆瓣电影 TOP250
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

	pages := parse.GetPages(BaseURL)
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
