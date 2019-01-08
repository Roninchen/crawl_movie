package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/parnurzeal/gorequest"
	"log"
	"strings"
	"net/url"
	"net/http"
	"time"
	"fmt"
	"math/rand"
)

// IP struct
type IP struct {
	ID   bson.ObjectId `bson:"_id" json:"-"`
	Data string        `bson:"data" json:"ip"`
	Type string        `bson:"type" json:"type"`
}

// NewIP .
func NewIP() *IP {
	return &IP{
		ID: bson.NewObjectId(),
	}
}
// IP66 get ip from 66ip.cn
func IP66() ([]string) {
	pollURL := "http://www.66ip.cn/mo.php?tqsl=1000"
	_, body, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return nil
	}
	body = strings.Split(body, "var mediav_ad_height = '60';")[1]
	body = strings.Split(body, "</script>")[1]
	body = strings.Split(body, "</div>")[0]
	body = strings.TrimSpace(body)
	body = strings.Replace(body, "	", "", -1)
	temp := strings.Split(body, "<br />")
	for index := 0; index < len(temp[:len(temp)-1]); index++ {
		AddIP(temp[index])
	}
	fmt.Print(temp)
	log.Println("IP66 done.")
	return temp
}
//validate and set ip to redis
//func Use66Ip() {
//	    temp := IP66()
//		//urli := url.URL{}
//		for index := 0; index < len(temp[:len(temp)-1]); index++{
//			//urlproxy, _ := urli.Parse(temp[index])
//			//c := http.Client{
//			//	Transport: &http.Transport{
//			//		Proxy: http.ProxyURL(urlproxy),
//			//	},
//			//}
//			//if resp, err := c.Get("https://movie.douban.com/subject/27133303/"); err != nil {
//			//	log.Fatalln(err)
//			//} else {
//				AddIP(temp[index])
//				//defer resp.Body.Close()
//				//body, _ := ioutil.ReadAll(resp.Body)
//				//fmt.Printf("%s\n", body)
//			//}
//		}
//	}


/**
* 返回response
*/
func GetRep(urls string) *http.Response {
	ip := string(GETIP())
	fmt.Print(ip)
	request, _ := http.NewRequest("GET", urls, nil)
	//随机返回User-Agent 信息
	request.Header.Set("User-Agent", getAgent())
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	proxy, err := url.Parse(ip)
	//设置超时时间
	timeout := time.Duration(20* time.Second)
	fmt.Printf("使用代理:%s\n",proxy)
	client := &http.Client{}
	if ip != "local"{
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
			Timeout: timeout,
		}
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200{

		fmt.Printf("line-99:遇到了错误-并切换ip %s\n",err)
		length := GetIPQueueLength()
		if length == 0{
			fmt.Printf("ip池暂时为空")
			time.Sleep(1000*time.Millisecond)
		}
		GETIP()
		GetRep(urls)
		//getIp(returnIP())

	}

	return response
}


/**
* 随机返回一个User-Agent
*/
func getAgent() string {
	agent  := [...]string{
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