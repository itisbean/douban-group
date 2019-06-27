package agent

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// GetHTML 获取html
func GetHTML(baseURL string) *http.Response {

	proxyurl := "http://122.136.212.132:53281" //代理IP，需要自己更换

	proxy, _ := url.Parse(proxyurl) // 解析代理IP

	netTransport := &http.Transport{ //要管理代理、TLS配置、keep-alive、压缩和其他设置，可以创建一个Transport
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 2, //超时设置
	}

	client := &http.Client{ //要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client
		Timeout:   time.Second * 2,
		Transport: netTransport,
	}
	req, err := http.NewRequest("GET", baseURL, nil) //NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。

	if err != nil {
		log.Println(err)
	}

	req.Header.Add("User-Agent", getAgent()) //模拟浏览器User-Agent
	resp, err := client.Do(req)              //Do方法发送请求，返回HTTP回复
	if err != nil {
		log.Println(err)
	}

	return resp //返回网页响应

}

func getAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}
