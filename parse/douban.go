package parse

import (
	"log"
	"math"
	"time"

	//"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DoubanGroupDbhyz 豆瓣小组
type DoubanGroupDbhyz struct {
	ID           uint `gorm:"primary_key"`     
	TopicID      int
	Topic        string
	AuthorID     int
	Author       string
	CreateTime   time.Time `gorm:"default:null"`
	NewReplyTime string 
	Reply        int
	Liked        int `gorm:"default:0"`
	Collect      int `gorm:"default:0"`
	Sharing      int `gorm:"default:0"`
	URL          string
	Content      string `gorm:"default:null"`
	Version      int
	IsDel        int `gorm:"default:0"`
}

// Page 分页
type Page struct {
	Page int
	URL  string
}

// GetTotal 获取总页数
func GetTotal(url string) int {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	totalstr, _ := doc.Find("#content > div > div.article > div.paginator > span.thispage").Attr("data-total-page")
	total, _ := strconv.Atoi(totalstr)

	return total
}

// Pages 分析分页
func Pages(url string, total int) (pages [][]Page) {
	size := 25 //每页25条，每25页一组

	lastKey := 0
	var pageList []Page

	for i := 0; i < total; i++ {
		key := int(math.Floor(float64(i / size)))
		if key != lastKey {
			pages = append(pages, pageList)
			pageList = append([]Page{})
			lastKey = key
		}
		pageList = append(pageList, Page{
			Page: i + 1,
			URL:  url + "?start=" + strconv.Itoa((total-i-1)*size),
		})
	}

	pages = append(pages, pageList)

	return pages
}

// PagesAll 获取全部的，包括漏页
func PagesAll(version []int) (pages [][]int) {
	size := 25 //每页25条，每25页一组

	lastKey := 0
	var pageList []int

	for i, v := range version {
		key := int(math.Floor(float64(i / size)))
		if key != lastKey {
			pages = append(pages, pageList)
			pageList = append([]int{})
			lastKey = key
		}
		pageList = append(pageList, v)
	}

	pages = append(pages, pageList)

	return pages
}

// ContentAll 分组
func ContentAll(items []DoubanGroupDbhyz) (groupItems [][]DoubanGroupDbhyz) {
	size := 30 //每30条一组

	lastKey := 0
	var itemList []DoubanGroupDbhyz

	for i, v := range items {
		key := int(math.Floor(float64(i / size)))
		if key != lastKey {
			groupItems = append(groupItems, itemList)
			itemList = append([]DoubanGroupDbhyz{})
			lastKey = key
		}
		itemList = append(itemList, v)
	}

	groupItems = append(groupItems, itemList)

	return groupItems
}

// Topics 分析话题
func Topics(doc *goquery.Document, version int) (items []DoubanGroupDbhyz, newVersion int) {
	//当前总页数
	totalstr, _ := doc.Find("#content > div > div.article > div.paginator > span.thispage").Attr("data-total-page")
	newVersion, _ = strconv.Atoi(totalstr)

	doc.Find("#content > div > div.article > div > table.olt > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 { //i=0时是标题列
			topicurl, _ := s.Find("td a").Eq(0).Attr("href")
			topicstr := strings.Split(topicurl, "/topic/")[1]
			topicstr = strings.Replace(topicstr, "/", "", -1)
			topicid, _ := strconv.Atoi(topicstr)

			topic, _ := s.Find("td a").Eq(0).Attr("title")

			authorurl, _ := s.Find("td a").Eq(1).Attr("href")
			authorstr := strings.Split(authorurl, "/people/")[1]
			authorstr = strings.Replace(authorstr, "/", "", -1)
			authorid, _ := strconv.Atoi(authorstr)

			author := s.Find("td a").Eq(1).Text()

			reply, _ := strconv.Atoi(s.Find("td").Eq(2).Text())

			timestr := s.Find("td").Eq(3).Text()
			if strings.Count(timestr, "-") == 1 {
				year := strconv.Itoa(time.Now().Year())
				timestr = strings.Join([]string{year, timestr}, "-")
			}
			if strings.Count(timestr, ":") == 0 {
				timestr = strings.Join([]string{timestr, "00:00:00"}, " ")
			}
			//newreplytime, _ := time.ParseInLocation("2006-01-02 15:04:05", timestr, time.Local)

			item := DoubanGroupDbhyz{
				TopicID:      topicid,
				Topic:        topic,
				AuthorID:     authorid,
				Author:       author,
				NewReplyTime: timestr,
				Reply:        reply,
				URL:          topicurl,
				Version:      version,
			}

			//log.Printf("i: %d, item: %v", i, item)
			items = append(items, item)
		}
	})

	return items, newVersion
}

// Detail 详情页
func Detail(doc *goquery.Document, item DoubanGroupDbhyz) DoubanGroupDbhyz {
	//item.URL = "https://www.douban.com/group/topic/143489532/"

	delText := doc.Find("#wrapper > div > ul > li").Eq(0).Find("p").Text()
	if delText == "呃...你想要的东西不在这儿" {
		item.IsDel = 1;
		return item
	}

	topicContent := doc.Find("#content > div > div.article > div.topic-content")

	timestr := topicContent.Find("div.topic-doc > h3 > span.color-green").Text()
	createtime, _ := time.ParseInLocation("2006-01-02 15:04:05", timestr, time.Local)

	mainContent := topicContent.Find("div.topic-doc > div#link-report > div.topic-content")

	images := ""
	mainContent.Find("div.topic-richtext > div.image-container > div.image-wrapper").Each(func(i int, s *goquery.Selection) {
		imgurl, _ := s.Find("img").Eq(0).Attr("src")
		if i == 0 {
			images = imgurl
		} else {
			images += "," + imgurl
		}
	})

	content := strings.TrimSpace(mainContent.Text())
	if images != "" {
		content = "[images]" + images + ";" +content
	}

	// TODO 点赞、收藏、转发 需要登录才能获取
	// liked, _ := strconv.Atoi(topicContent.Find("div.sns-bar > div.action-react > a > span.react-num").Text())
	// collect, _ := strconv.Atoi(topicContent.Find("div.sns-bar > div.action-collect > a > span.react-num").Text())
	// sharing, _ := strconv.Atoi(topicContent.Find("div.sns-bar > div.sharing > div > div > div > span > a > span.rec-num").Text())  

	item.CreateTime = createtime
	item.Content = content

	return item
}

