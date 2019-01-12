package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) Concurrency() {
	url :="http://localhost:8080/crawl_movie"
	i, err := c.GetInt(":count")
	if err!=nil {
		logs.Error(err)
	}
	logs.Info(i)
	for j:=0; j<i;j++  {
		logs.Info("第%d个启动",j)
		go http.Get(url)
	}

}
