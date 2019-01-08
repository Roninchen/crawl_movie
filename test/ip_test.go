package test

import (
	"testing"
	"crawl_movie/models"
	"github.com/astaxie/beego/logs"
	"time"
	"fmt"
)

func TestIp(t *testing.T)  {
	models.ConnectRedis("127.0.0.1:6379")
	go models.IP66()
	sUrl := "https://movie.douban.com/subject/27133303/"
	models.PutinQueue(sUrl)
	length := models.GetQueueLength()
	if length == 0 {
		fmt.Print("kong")
	}

	sUrl = models.PopfromQueue()
	//我们应当判断sUrl是否应该被访问过
	if models.IsVisit(sUrl) {
		fmt.Print("wu")
	}
	logs.SetLogger(sUrl)
	go models.Run(sUrl)
	time.Sleep(2000)
}
