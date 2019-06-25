// 爬取豆瓣电影 TOP250
package main

import (
	"strconv"
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

// Add 新增数据
func Add(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		if err := model.DB.Create(&topic).Error; err != nil {
			log.Printf("db.Create index: %s, err : %v", strconv.Itoa(index), err)
		}
	}
}

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

	Add(topics)
}

func main() {
	Start()

	defer model.DB.Close()
}
