package controllers

import (
	"crawl_movie/models"
	"github.com/spf13/viper"
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
	models.ConnectRedis(viper.GetString("redis.host")+":"+viper.GetString("redis.port"))

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
		if !URLFilter(sUrl){
			continue
		}
		if strings.Contains(sUrl, "subject") {
			logs.Info(sUrl)
			go models.Run(sUrl)
			time.Sleep(time.Second*3)
		}
	}
}

func URLFilter(sUrl string) bool {
	if strings.Contains(sUrl,"celebrity") {
		logs.Info("contains celebrity continue!")
		return false
	}
	if strings.Contains(sUrl,"short_video"){
		logs.Info("contains short_video continue!")
		return false
	}
	if strings.Contains(sUrl, "photos") {
		logs.Info("contains photos continue!")
		return false
	}
	if strings.Contains(sUrl, "mupload") {
		logs.Info("contains mupload continue!")
		return false
	}
	if strings.Contains(sUrl, "trailer") {
		logs.Info("contains trailer continue!")
		return false
	}
	if strings.Contains(sUrl, "video") {
		logs.Info("contains video continue!")
		return false
	}
	if strings.Contains(sUrl, "review") {
		logs.Info("contains review continue!")
		return false
	}
	if strings.Contains(sUrl, "questions") {
		logs.Info("contains questions continue!")
		return false
	}
	if strings.Contains(sUrl, "discussion") {
		logs.Info("contains discussion continue!")
		return false
	}
	if strings.Contains(sUrl, "doulists") {
		logs.Info("contains doulists continue!")
		return false
	}
	if strings.Contains(sUrl, "awards") {
		logs.Info("contains awards continue!")
		return false
	}
	if strings.Contains(sUrl, "wishes") {
		logs.Info("contains awards continue!")
		return false
	}
	if strings.Contains(sUrl, "collections") {
		logs.Info("contains collections continue!")
		return false
	}
	if strings.Contains(sUrl, "update_image") {
		logs.Info("contains update_image continue!")
		return false
	}
	if strings.Contains(sUrl, "comments") {
		logs.Info("contains comments continue!")
		return false
	}
	return true
}