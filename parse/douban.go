package parse

import (
	"log"
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
	Forward      int
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
	pages = append(pages, Page{Page: 1, URL: url})
	doc.Find("#content > div > div.article > div.paginator > a").Each(func(i int, s *goquery.Selection) {
		if i < 2 {
			page, _ := strconv.Atoi(s.Text())
			url, _ := s.Attr("href")

			pages = append(pages, Page{
				Page: page,
				URL:  url,
			})
		}
	})
	return pages
}

// Topics 分析话题
func Topics(doc *goquery.Document) (items []DoubanGroupDbhyz) {
	doc.Find("#content > div > div.article > div > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			topicstr, _ := s.Find("td a").Eq(0).Attr("href")
			topicstr = strings.TrimLeft(topicstr, "/topic/")
			topicid, _ := strconv.Atoi(topicstr)

			topic, _ := s.Find("td a").Eq(0).Attr("title")

			authorstr, _ := s.Find("td a").Eq(1).Attr("href")
			authorstr = strings.TrimLeft(topicstr, "/people/")
			authorid, _ := strconv.Atoi(authorstr)

			author := s.Find("td a").Eq(1).Text()

			reply, _ := strconv.Atoi(s.Find("td").Eq(2).Text())

			newreplytime := s.Find("td").Eq(3).Text()

			//other = strings.TrimLeft(other, "  / ")

			// desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
			// DescInfo := strings.Split(desc, "\n")
			// desc = DescInfo[0]

			// movieDesc := strings.Split(DescInfo[1], "/")
			// year := strings.TrimSpace(movieDesc[0])
			// area := strings.TrimSpace(movieDesc[1])
			// tag := strings.TrimSpace(movieDesc[2])

			// star := s.Find(".bd .star .rating_num").Text()

			// comment := strings.TrimSpace(s.Find(".bd .star span").Eq(3).Text())
			// compile := regexp.MustCompile("[0-9]")
			// comment = strings.Join(compile.FindAllString(comment, -1), "")

			item := DoubanGroupDbhyz{
				TopicID: 	  topicid,
				Topic:        topic,
				AuthorID: 	  authorid,
				Author: 	  author,
				CreateTime:   "",
				NewReplyTime: newreplytime,
				Reply:        reply,
				Liked: 		  0,
				Collect: 	  0,
				Forward: 	  0,
				Content: 	  "",
			}

			log.Printf("i: %d, item: %v", i, item)
			items = append(items, item)
		}
	})

	return items
}
