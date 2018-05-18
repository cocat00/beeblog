package controllers

import "github.com/astaxie/beego"

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.Data["IsCategory"] = true

	c.TplName = "category.html"
}

func (c *CategoryController) Post()  {

}
