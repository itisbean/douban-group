package agent

import (
	//"math/rand"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"

	browser "github.com/EDDYCJY/fake-useragent"
	log "github.com/cihub/seelog"
)

// ProxyIP 代理Ip
type ProxyIP struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var ipPool []string

// GetHTML 获取html
func GetHTML(baseURL string, userAgent string, proxyAddr string) (*http.Response, error) {
	time.Sleep(time.Second * 1)

	proxy, _ := url.Parse(proxyAddr) // 解析代理IP

	netTransport := &http.Transport{ //要管理代理、TLS配置、keep-alive、压缩和其他设置，可以创建一个Transport
		Proxy:        http.ProxyURL(proxy),
		MaxIdleConns: 0,
		MaxIdleConnsPerHost:   0,
		ResponseHeaderTimeout: time.Second * 15, //超时设置
		IdleConnTimeout: 0,
		DisableKeepAlives: true,
	}

	client := &http.Client{ //要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client
		Timeout:   time.Second * 30,
		Transport: netTransport,
	}

	req, err := http.NewRequest("GET", baseURL, nil) //NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。
	if err != nil {
		return nil, err
	}

	if userAgent != "" {
		req.Header.Add("User-Agent", userAgent) //模拟浏览器User-Agent
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := client.Do(req) //Do方法发送请求，返回HTTP回复
	if err != nil {
		return nil, err
	}

	return resp, err //返回网页响应
}

// GetAgent 随机获取user-agent
func GetAgent() string {
	random := browser.Chrome()
	//log.Debugf("Random: %s", random)
	return random
}

// GetProxy 获取代理Ip
func GetProxy() (proxyurl string, useragent string) {
	localpool := "http://localhost:8090/get"

	try := 0

	for {
		if try >= 15 {
			break
		}
		try++

		proxyres, err := http.Get(localpool)
		if err != nil {
			log.Error(err)
			continue
		}

		defer proxyres.Body.Close()
		body, err := ioutil.ReadAll(proxyres.Body)
		if err != nil {
			log.Error(err)
			continue
		}

		proxyip := &ProxyIP{}
		json.Unmarshal(body, &proxyip)

		proxyurl = "http://" + proxyip.IP + ":" + strconv.Itoa(proxyip.Port)

		var ip, useragent = ProxyThorn(proxyurl)
		//判断是否有返回ip，并且请求状态为200
		if ip != "" && useragent != "" {
			log.Debugf(proxyip.IP + " 请求 http://icanhazip.com 返回ip:【" + ip + "】-【检测结果：可用】")
			ipPool = append(ipPool, proxyip.IP)
			break
		} else {
			proxyurl = ""
			log.Debugf(proxyip.IP + "【检测结果：不可用】")
		}
	}

	log.Info(ipPool)
	return proxyurl, useragent
}

// ProxyThorn 验证代理ip是否可用
// 通过传入一个代理ip，然后使用它去访问一个url看看是否访问成功，以此为依据进行判断当前代理ip是否有效。
// 参数：proxy_addr 要验证的ip
// 返回：ip 验证通过的ip、status 状态（200表示成功）
func ProxyThorn(proxyAddr string) (ip string, useragent string) {
	//访问查看ip的一个网址
	httpURL1 := "http://www.douban.com/group/639264/discussion"
	httpURL2 := "http://icanhazip.com"

	proxy, err := url.Parse(proxyAddr)

	netTransport := &http.Transport{
		Proxy:        http.ProxyURL(proxy),
		MaxIdleConns: 0,
		MaxIdleConnsPerHost:   0,
		ResponseHeaderTimeout: time.Second * 15,
		DisableKeepAlives: true,
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 15,
		Transport: netTransport,
	}

	res, err := httpClient.Get(httpURL2)
	if err != nil {
		log.Debug("检测失败，错误信息：", err)
		return "", ""
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Debugf("检测失败，http code: %d", res.StatusCode)
		return "", ""
	}

	useragent = GetAgent()
	resp, err := GetHTML(httpURL1, useragent, proxyAddr)
	if err != nil {
		log.Debugf("检测失败：douban test failed,", err)
		return "", ""
	}
	
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		log.Debugf("检测失败，http code:%d", resp.StatusCode)
		return "", ""
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()
	if err != nil {
		log.Debugf("检测失败，goquery读取doc失败：", err)
		return "", ""
	}
	length := doc.Find("#content > div > div.article > div > table.olt > tbody > tr").Length()
	if (length == 0) {
		log.Debugf("检测失败，无法获取页面可用内容")
		return "", ""
	}

	c, _ := ioutil.ReadAll(res.Body)
	return string(c), useragent
}
