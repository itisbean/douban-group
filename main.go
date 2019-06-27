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

	version, err := model.GetVersion()
	if err != nil {
		log.Println(err)
		return
	}

	total := parse.GetTotal(BaseURL, version)

	var pages [][]parse.Page
	pages = parse.Pages(BaseURL, (total - version + 1))

	for _, pageList := range pages {
		//获取新的Ip抓取页面；延时

		for _, page := range pageList {
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
}

func main() {
	Start()

	defer model.DB.Close()
}
