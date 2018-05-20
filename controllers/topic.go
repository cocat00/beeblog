package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get()  {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["IsTopic"] = true
	c.TplName = "topic.html"

	var err error
	c.Data["Topics"], err = models.GetAllTopics(false)

	if err != nil {
		beego.Error(err)
	}
}

func (c *TopicController) Post()  {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}


	// 解析表单
	tid := c.Input().Get("tid")
	title := c.Input().Get("title")
	content := c.Input().Get("content")

	var err error
	if len(tid) == 0 {
		err = models.AddTopic(title, content)
	} else {
		err = models.ModifyTopic(tid, title, content)
	}

	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic", 302)

}

func (c *TopicController) Add()  {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "topic_add.html"
}

func (c *TopicController) View()  {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "topic_view.html"

	tid := c.Ctx.Input.Param("0")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
}

func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"

	tid := this.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
	this.Data["Tid"] = tid
}


func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic", 302)
}