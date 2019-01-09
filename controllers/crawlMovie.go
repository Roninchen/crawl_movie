package controllers

import (
	"crawl_movie/models"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
)

type CrawlMovieController struct {
	beego.Controller
}

func (c *CrawlMovieController) CrawlMovie() {

	//连接到redis
	models.ConnectRedis("111.231.84.238:6379")

	//爬虫入口url
	sUrl := "https://movie.douban.com/subject/3878007/"
	models.PutinQueue(sUrl)
	// go models.IP66()
	for {
		//models.IP66()
		length := models.GetQueueLength()
		if length == 0 {
			break //如果url队列为空 则退出当前循环
		}

		sUrl = models.PopfromQueue()
		//我们应当判断sUrl是否应该被访问过
		if models.IsVisit(sUrl) {
			continue
		}
		if strings.Contains(sUrl,"celebrity") {
			logs.Info("contains celebrity continue!")
			continue
		}
		if strings.Contains(sUrl,"short_video"){
			logs.Info("contains short_video continue!")
			continue
		}
		if strings.Contains(sUrl, "photos") {
			logs.Info("contains photos continue!")
			continue
		}
		if strings.Contains(sUrl, "mupload") {
			logs.Info("contains mupload continue!")
			continue
		}
		if strings.Contains(sUrl, "trailer") {
			logs.Info("contains trailer continue!")
			continue
		}
		if strings.Contains(sUrl, "video") {
			logs.Info("contains video continue!")
			continue
		}
		if strings.Contains(sUrl, "review") {
			logs.Info("contains review continue!")
			continue
		}
		if strings.Contains(sUrl, "subject") {
			logs.Info(sUrl)
			go models.Run(sUrl)
			time.Sleep(time.Second*3)
		}
	}
}