package controllers

import (
	"crawl_movie/models"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
)

type CrawlMovieController struct {
	beego.Controller
}
type UpdateMovieController struct {
	beego.Controller
}
func (c *CrawlMovieController) CrawlMovie() {

	//连接到redis
	models.ConnectRedis(viper.GetString("redis.host")+":"+viper.GetString("redis.port"))


	//爬虫入口url
	sUrl := "https://movie.douban.com/subject/26266893/?from=playing_poster"
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
		if models.GetQueueLength()>100000 && !URLFilter(sUrl){
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

func(c *UpdateMovieController)UpdateMovie(){
	movieInfo, err := models.GetMovieGradeIsZero()
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": 1, "message": "err","data":err}
		c.ServeJSON()
		return
	}
	if len(movieInfo) < 1 {
		c.Data["json"] = map[string]interface{}{"success": 1, "message": "IS NULL","data":""}
		c.ServeJSON()
		return
	}
	if movieInfo == nil {
		c.Data["json"] = map[string]interface{}{"success": 1, "message": "IS NULL","data":""}
		c.ServeJSON()
		return
	}
	logs.Info("信息长度:",len(movieInfo))
	for _,v := range movieInfo{
		movieId := strconv.FormatInt(v.Movie_id,10)
		url := "https://movie.douban.com/subject/"+movieId +"/"
		movieInfo, err, _ := models.ReturnMovieInfoByUrl(url)
		if err != nil {
			logs.Info(err)
			//c.Data["json"] = map[string]interface{}{"success": 1, "message": "err","data":err}
			//c.ServeJSON()
			continue
		}
		if movieInfo == nil {
			logs.Info("movieInfo is null")
			//c.Data["json"] = map[string]interface{}{"success": 1, "message": "movieInfo is null","data":err}
			//c.ServeJSON()
			continue
		}
		movieInfo.Id = v.Id
		err = models.UpdateMovieGrade(movieInfo)
		if err != nil {
			logs.Info(err)
			//c.Data["json"] = map[string]interface{}{"success": 1, "message": "err","data":err}
			//c.ServeJSON()
			continue
		}
		//c.Data["json"] = map[string]interface{}{"success": 0, "message": "ok","data":movieInfo}
		//c.ServeJSON()
		//return
	}
	c.Data["json"] = map[string]interface{}{"success": 0, "message": "ok","data":movieInfo}
	c.ServeJSON()
	return
}