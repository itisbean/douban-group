package parse

import (
	"log"
	"time"

	//"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DoubanGroupDbhyz 豆瓣小组
type DoubanGroupDbhyz struct {
	TopicID      int
	Topic        string
	AuthorID     int
	Author       string
	CreateTime   string
	NewReplyTime string
	Reply        int
	Liked        int
	Collect      int
	Sharing      int
	URL          string
	Content      string
}

// Page 分页
type Page struct {
	Page int
	URL  string
}

// GetPages 获取分页
func GetPages(url string) []Page {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	return Pages(doc, url)
}

// Pages 分析分页
func Pages(doc *goquery.Document, url string) (pages []Page) {
	size := 25

	totalstr, _ := doc.Find("#content > div > div.article > div.paginator > span.thispage").Attr("data-total-page")
	total, _ := strconv.Atoi(totalstr)

	for i := 1; i <= total; i++ {
		pages = append(pages, Page{
			Page: i,
			URL:  url + "?start=" + strconv.Itoa((i-1)*size),
		})
	}

	return pages
}

// Topics 分析话题
func Topics(doc *goquery.Document) (items []DoubanGroupDbhyz) {
	doc.Find("#content > div > div.article > div > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
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

			newreplytime := s.Find("td").Eq(3).Text()
			if strings.Count(newreplytime, "-") == 1 {
				year := strconv.Itoa(time.Now().Year())
				newreplytime = strings.Join([]string{year, newreplytime}, "-")
			}

			item := DoubanGroupDbhyz{
				TopicID:      topicid,
				Topic:        topic,
				AuthorID:     authorid,
				Author:       author,
				CreateTime:   "",
				NewReplyTime: newreplytime,
				Reply:        reply,
				Liked:        0,
				Collect:      0,
				Sharing:      0,
				URL:          topicurl,
				Content:      "",
			}

			item = Detail(item)

			log.Printf("i: %d, item: %v", i, item)
			items = append(items, item)
		}
	})

	return items
}

// Detail 详情页
func Detail(item DoubanGroupDbhyz) DoubanGroupDbhyz {
	//item.URL = "https://www.douban.com/group/topic/143489532/"
	doc, err := goquery.NewDocument(item.URL)
	if err != nil {
		log.Println(err)
	}

	topicContent := doc.Find("#content > div > div.article > div.topic-content")

	createtime := topicContent.Find("div.topic-doc > h3 > span.color-green").Text()
	content := strings.TrimSpace(topicContent.Find("div.topic-doc > div#link-report > div.topic-content").Text())

	// TODO 点赞、收藏、转发 需要登录才能获取
	// liked, _ := strconv.Atoi(topicContent.Find("div.sns-bar > div.action-react > a > span.react-num").Text())
	// collect, _ := strconv.Atoi(topicContent.Find("div.sns-bar > div.action-collect > a > span.react-num").Text())
	// sharing, _ := strconv.Atoi(topicContent.Find("div.sns-bar > div.sharing > div > div > div > span > a > span.rec-num").Text())

	item.CreateTime = createtime
	// item.Liked = liked
	// item.Collect = collect
	// item.Sharing = sharing
	item.Content = content

	return item
}
