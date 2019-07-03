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

	newVersion := parse.GetTotal(BaseURL)

	version := model.GetVersion(newVersion)
	
	var pages [][]parse.Page
	pages = parse.PagesAll(BaseURL, newVersion, version)

	log.Printf("pages group:%d", len(pages))

	for _, pageList := range pages {
		wg.Add(len(pageList))
		go func(pageList []parse.Page) {

			//1、获取新的Ip和user-agent抓取页面；延时防封禁；
			proxyAddr, userAgent := agent.GetProxy() //代理IP，需要自己更换
			if proxyAddr == "" {
				log.Println("无法获取代理Ip，请稍后重试")
				return
			}

			var items []parse.DoubanGroupDbhyz
			//2、开始抓取每页话题
			for _, page := range pageList {
				defer wg.Done()

				resp := agent.GetHTML(page.URL, userAgent, proxyAddr)
				if resp == nil {
					log.Println("Get Html Error,Please Retry")
					return
				}

				log.Printf("http code:%d", resp.StatusCode)

				if resp.StatusCode == 403 {
					log.Println("403 Forbidden,Please Retry")
					return
				}
				doc, err := goquery.NewDocumentFromResponse(resp)
				defer resp.Body.Close()

				if err != nil {
					log.Println(err)
					return
				}

				curVersion := newVersion - page.Page + 1
				//items = append(items, parse.Topics(doc, curVersion)...)
				items = parse.Topics(doc, curVersion)
				log.Printf("items:%v", items)
				model.Save(items)
			}
		}(pageList)

		time.Sleep(time.Second * 5)
	}

	wg.Wait()
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
