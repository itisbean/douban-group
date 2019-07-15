// 爬取豆瓣小组
package main

import (
	"time"
	"go-crawler/douban-group/agent"
	"strconv"

	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"go-crawler/douban-group/model"
	"go-crawler/douban-group/parse"

	log "github.com/cihub/seelog"
)

// 抓取网站：豆瓣🔥小组
var (
	BaseURL    = "https://www.douban.com/group/639264/discussion"
	newVersion = 0

	wg sync.WaitGroup
)

func curVersion() (v []int) {
	try := 9
	size := 900
	v = model.GetVersion(size*(try-1), size*(try-1)+size)
	// for i := 0; i <= size; i++ {
	// 	v = append(v, ((try-1)*size + i))
	// }
	return v
}

// Start1 分页抓取帖子（ID、标题、作者、最后回复时间等）
func Start1() {

	version := curVersion()

	log.Debug(version)
	//return

	if len(version) == 0 {
		return
	}

	var pages [][]int
	pages = parse.PagesAll(BaseURL, version)

	log.Info("pages group:", len(pages))

	newVersion = parse.GetTotal(BaseURL)

	for _, pageList := range pages {
		wg.Add(len(pageList))
		go func(pageList []int) {

			//1、获取新的Ip和user-agent抓取页面；延时防封禁；
			proxyAddr, userAgent := agent.GetProxy() //代理IP，需要自己更换
			if proxyAddr == "" {
				log.Error("无可用代理Ip，请稍后重试")
				log.Info("failed:", pageList)
				defer wg.Add(-len(pageList))
				return
			}

			//var items []parse.DoubanGroupDbhyz
			var failed []int
			//2、开始抓取每页话题
			for _, page := range pageList {
				defer wg.Done()

				log.Debug("total:", newVersion)
				curURL := BaseURL + "?start=" + strconv.Itoa((newVersion-page)*25)

				resp, err := agent.GetHTML(curURL, userAgent, proxyAddr)
				if resp == nil {
					failed = append(failed, page)
					log.Error("Get请求页面失败，", err)
					continue
				}

				log.Debug("http code:", resp.StatusCode)

				if resp.StatusCode == 403 {
					failed = append(failed, page)
					log.Error("错误403 Forbidden,请更换Ip")
					continue
				}
				doc, err := goquery.NewDocumentFromResponse(resp)
				defer resp.Body.Close()

				if err != nil {
					failed = append(failed, page)
					log.Critical(err)
					continue
				}

				items, total := parse.Topics(doc, page)
				if len(items) == 0 {
					failed = append(failed, page)
					log.Error("爬虫解析失败，内容返回为空")
					continue
				}
				if total > newVersion {
					newVersion = total
				}
				log.Debug("items:", items)
				log.Debug("new version:", newVersion)
				model.Save(items)
			}

			log.Info("failed:", failed)
		}(pageList)

		//time.Sleep(time.Second * 5)
	}

	wg.Wait()
}

// Start2 从数据库获取未抓内容的话题，进入详情页抓取内容
func Start2() {
	//1、获取当前version下的话题数据

	//2、按30条一组分割

	//3、循环抓取，每组更新一次ip和设置延时，保存数据
}

// SetLogger 初始化日志配置
func SetLogger(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		logger, err := log.LoggerFromConfigAsFile(fileName)
		if err != nil {
			panic(err)
		}

		log.ReplaceLogger(logger)
	}
	log.Info("log initialize finish.")
}

func main() {
	bT := time.Now()

	SetLogger("logConfig.xml")
	defer log.Flush()

	Start1()
	//Start2()

	eT := time.Since(bT)

	log.Info("run time:", eT)

	defer model.DB.Close()
}
