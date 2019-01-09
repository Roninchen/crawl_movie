package models

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2/bson"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"crawl_movie/util"
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
type Result struct {
	Ip       string  `json:"ip"`
	Port     int     `json:"port"`
	Location string  `json:"location,omitempty"`
	Source   string  `json:"source"`
	Speed    float64 `json:"speed,omitempty"`
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



/**
* 返回response
*/
func GetRep(urls string) *http.Response {
	ip :=ReturnIp()
	logs.Info("使用代理:%s\n",ip)
	resp, _, errs := gorequest.New().
		Proxy(ip).
		Get(urls).
		Set("User-Agent", util.RandomUA()).
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8").
		Timeout(time.Second * 6).
		End()

	if errs != nil || resp.StatusCode != 200{

		logs.Info("line-99:遇到了错误-并切换ip %s\n",errs)
		if strings.Contains(urls, "subject") {
			return GetRep(urls)
		}
		return nil
	}

	return resp
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

func ReturnIp() string {
	var result = new(Result)
	ipurl := "http://localhost:8090/get"
	resp, err := http.Get(ipurl)
	if err != nil{
		fmt.Println(err)
	}
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	logs.Info(body)
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		fmt.Println(err)
	}
	ip :=  "http://" + result.Ip + ":" + strconv.Itoa(result.Port)
	return ip
}