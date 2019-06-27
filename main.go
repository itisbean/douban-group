// 爬取豆瓣小组
package main

import (
	"log"
	"sync"

	//"strings"

	"github.com/PuerkitoBio/goquery"

	"go-crawler/douban-group/model"
	"go-crawler/douban-group/parse"
)

// 抓取网站：豆瓣🔥小组
var (
	BaseURL = "https://www.douban.com/group/639264/discussion"

	wg sync.WaitGroup
)

// Start 开始爬取
func Start() {
	//var topics []parse.DoubanGroupDbhyz

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
		wg.Add(1)
		go func(page parse.Page) {
			defer wg.Done()

			doc, err := goquery.NewDocument(page.URL)
			if err != nil {
				log.Println(err)
			}

			model.Save(parse.Topics(doc))
		}(page)

		wg.Wait()
	}
}

func main() {
	Start()

	defer model.DB.Close()
}
