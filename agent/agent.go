package agent

import (
	"github.com/PuerkitoBio/goquery"
	//"math/rand"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
)

// ProxyIP 代理Ip
type ProxyIP struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// GetHTML 获取html
func GetHTML(baseURL string, userAgent string, proxyAddr string) *http.Response {
	time.Sleep(time.Second * 2)

	proxy, _ := url.Parse(proxyAddr) // 解析代理IP

	netTransport := &http.Transport{ //要管理代理、TLS配置、keep-alive、压缩和其他设置，可以创建一个Transport
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 10, //超时设置
	}

	client := &http.Client{ //要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, err := http.NewRequest("GET", baseURL, nil) //NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。
	if err != nil {
		log.Println("【GetHtml】错误信息1：", err)
	}

	if (userAgent != "") {
		req.Header.Add("User-Agent", userAgent) //模拟浏览器User-Agent
	}
	
	resp, err := client.Do(req)             //Do方法发送请求，返回HTTP回复
	if err != nil {
		log.Println("【GetHtml】错误信息2：", err)
	}

	return resp //返回网页响应
}

// GetAgent 随机获取user-agent
func GetAgent() string {
	// var uas = [...]string{
	// 	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36",
	// 	//"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.76 Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	// 	"Mozilla/5.0 (Linux; Android 7.0; SM-G570M Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Mobile Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36",
	// 	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586",
	// 	"Mozilla/5.0 (Linux; Android 6.0; MYA-L22 Build/HUAWEIMYA-L22) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36",
	// 	//"Mozilla/5.0 (Linux; Android 6.0; vivo 1713 Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.124 Mobile Safari/537.36",
	// 	//"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	// }
	// n := rand.Intn(len(uas))
	//random := uas[n]
	random := browser.Chrome()
	log.Printf("Random: %s", random)
	return random
}

// GetProxy 获取代理Ip
func GetProxy() (proxyurl string, useragent string) {
	localpool := "http://localhost:8090/get"

	try := 0

	for {
		if try >= 10 {
			break
		}
		try++

		proxyres, err := http.Get(localpool)
		if err != nil {
			log.Println(err)
			continue
		}

		defer proxyres.Body.Close()
		body, err := ioutil.ReadAll(proxyres.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		proxyip := &ProxyIP{}
		json.Unmarshal(body, &proxyip)

		proxyurl = "http://" + proxyip.IP + ":" + strconv.Itoa(proxyip.Port)

		var ip, status, useragent = ProxyThorn(proxyurl)
		//判断是否有返回ip，并且请求状态为200
		if status == 200 && ip != "" && useragent != "" {
			log.Println(proxyip.IP + " 请求 http://icanhazip.com 返回ip:【" + ip + "】-【检测结果：可用】")
			break
		} else {
			log.Println(proxyip.IP + "【检测结果：不可用】")
		}
	}

	return proxyurl, useragent
}

// ProxyThorn 验证代理ip是否可用
// 通过传入一个代理ip，然后使用它去访问一个url看看是否访问成功，以此为依据进行判断当前代理ip是否有效。
// 参数：proxy_addr 要验证的ip
// 返回：ip 验证通过的ip、status 状态（200表示成功）
func ProxyThorn(proxyAddr string) (ip string, status int, useragent string) {
	//访问查看ip的一个网址
	httpURL1 := "https://www.douban.com/group/639264/discussion"
	httpURL2 := "http://icanhazip.com"

	proxy, err := url.Parse(proxyAddr)

	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 5,
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 5,
		Transport: netTransport,
	}

	res, err := httpClient.Get(httpURL2)
	if err != nil {
		log.Println("检测失败，错误信息：",err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(err)
		return
	}

	useragent = GetAgent()
	resp := GetHTML(httpURL1, useragent, proxyAddr)
	if resp == nil {
		log.Printf("检测失败：douban test failed")
		return
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		log.Printf("检测失败，http code:%d", resp.StatusCode)
		return
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()
	if err != nil {
		log.Println("检测失败，goquery读取doc失败：", err)
		return
	}
	length := doc.Find("#content > div > div.article > div > table.olt > tbody > tr").Length()
	if (length == 0) {
		log.Println("检测失败，无法获取页面可用内容")
		return
	}

	c, _ := ioutil.ReadAll(res.Body)
	return string(c), res.StatusCode, useragent
}
