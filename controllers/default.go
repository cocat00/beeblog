package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["IsHome"] = true
	//c.Data["IsCategory"] = true
	c.TplName = "home.html"

	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var err error
	c.Data["Topics"], err = models.GetAllTopics(true)

	if err != nil {
		beego.Error(err)
	}
}
