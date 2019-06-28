package agent

import (
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
	time.Sleep(time.Second * 1)

	proxy, _ := url.Parse(proxyAddr) // 解析代理IP

	netTransport := &http.Transport{ //要管理代理、TLS配置、keep-alive、压缩和其他设置，可以创建一个Transport
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 20, //超时设置
	}

	client := &http.Client{ //要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client
		Timeout:   time.Second * 20,
		Transport: netTransport,
	}

	req, err := http.NewRequest("GET", baseURL, nil) //NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("User-Agent", userAgent) //模拟浏览器User-Agent
	resp, err := client.Do(req)             //Do方法发送请求，返回HTTP回复
	if err != nil {
		log.Println(err)
	}

	return resp //返回网页响应
}

// GetAgent 随机获取user-agent
func GetAgent() string {
	random := browser.Chrome()
	log.Printf("Random: %s", random)
	return random
}

// GetProxy 获取代理Ip
func GetProxy() string {
	localpool := "http://192.168.254.128:8090/get"

	proxyurl := ""
	try := 0

	for {
		if try >= 10 {
			break
		}

		proxyres, err := http.Get(localpool)
		if err != nil {
			log.Println(err)
		}

		defer proxyres.Body.Close()
		body, err := ioutil.ReadAll(proxyres.Body)
		if err != nil {
			log.Println(err)
		}

		proxyip := &ProxyIP{}
		json.Unmarshal(body, &proxyip)

		proxyurl = "http://" + proxyip.IP + ":" + strconv.Itoa(proxyip.Port)

		var ip, status = ProxyThorn(proxyurl)
		//判断是否有返回ip，并且请求状态为200
		if status == 200 && ip != "" {
			log.Println(proxyip.IP + " 请求 http://icanhazip.com 返回ip:【" + ip + "】-【检测结果：可用】")
			break
		} else {
			log.Println(proxyip.IP + "【检测结果：不可用】")
		}

		try++
	}

	return proxyurl
}

// ProxyThorn 验证代理ip是否可用
// 通过传入一个代理ip，然后使用它去访问一个url看看是否访问成功，以此为依据进行判断当前代理ip是否有效。
// 参数：proxy_addr 要验证的ip
// 返回：ip 验证通过的ip、status 状态（200表示成功）
func ProxyThorn(proxyAddr string) (ip string, status int) {
	//访问查看ip的一个网址
	httpURL := "http://icanhazip.com"
	proxy, err := url.Parse(proxyAddr)

	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	res, err := httpClient.Get(httpURL)
	if err != nil {
		//fmt.Println("错误信息：",err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(err)
		return
	}
	c, _ := ioutil.ReadAll(res.Body)
	return string(c), res.StatusCode
}
