package routers

import (
	"crawl_movie/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/crawl_movie", &controllers.CrawlMovieController{}, "*:CrawlMovie")
    beego.Router("/concurrency/:count",&controllers.MainController{},"get:Concurrency")
    beego.Router("/update",&controllers.UpdateMovieController{},"*:UpdateMovie")
    //beego.Router("")
}
