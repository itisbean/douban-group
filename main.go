// 爬取豆瓣小组
package main

import (
	"go-crawler/douban-group/agent"
	"log"
	"sync"
	"time"

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

// Start1 分页抓取帖子（ID、标题、作者、最后回复时间等）
func Start1() {
	version := model.GetVersion()

	newVersion := parse.GetTotal(BaseURL, version)
	newVersion = 1

	var pages [][]parse.Page
	pages = parse.Pages(BaseURL, (newVersion - version + 1))

	for _, pageList := range pages {
		//1、获取新的Ip和user-agent抓取页面；延时防封禁；
		proxyAddr := agent.GetProxy() //代理IP，需要自己更换
		userAgent := agent.GetAgent()

		//2、开始抓取每页话题
		for index, page := range pageList {
			wg.Add(1)
			go func(page parse.Page) {
				defer wg.Done()

				resp := agent.GetHTML(page.URL, userAgent, proxyAddr)
				if resp.StatusCode == 403 {
					log.Println("403 Forbidden,Please Retry")
					return
				}
				doc, err := goquery.NewDocumentFromResponse(resp)
				if err != nil {
					log.Println(err)
					return
				}

				model.Save(parse.Topics(doc, index))
			}(page)

			wg.Wait()
		}

		time.Sleep(time.Second * 5)
	}
}

// Start2 从数据库获取未抓内容的话题，进入详情页抓取内容
func Start2() {
	//1、获取当前version下的话题数据

	//2、按30条一组分割

	//3、循环抓取，每组更新一次ip和设置延时，保存数据
}

func main() {
	Start1()
	//Start2()

	defer model.DB.Close()
}
