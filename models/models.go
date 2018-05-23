package models

import (
	"time"
	"github.com/Unknwon/com"
	"os"
	"path"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3" //导入驱动
	"strconv"
)

const (
	_DB_NAME	    = "data/beeblog"
	_SQLITES_DRIVER = "sqlite3"
)

type Category struct {
	Id		 		int64
	Title    		string
	Created 		time.Time
	Views   		int64     	`orm:"index"`
	TopicTime       time.Time   `orm:"index"`
	TopicCount 		int64
	TopicLastUserId int64
}

type Topic struct {
	Id 				int64
	Uid 			int64
	Title 			string
	Category		string
	Content 		string      `orm:"size(5000)"`
	Attachment 		string
	Created		    time.Time   `orm:"index"`
	Updated 	 	time.Time   `orm:"index"`
	Views  			int64       `orm:"index"`
	Author 			string
	ReplyTime 		time.Time   `orm:"index"`
	ReplyCount 		int64
	ReplyLastUserId int64
}

func RegisterDB()  {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册定义的 model
	orm.RegisterModel(new(Category), new(Topic))
	//注册驱动
	orm.RegisterDriver(_SQLITES_DRIVER, orm.DRSqlite)
	// 设置默认数据库
	//数据库存放位置：_DB_NAME ， 数据库别名：default
	orm.RegisterDataBase("default", _SQLITES_DRIVER, _DB_NAME, 10)
}


func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{Title: name, Created:time.Now(), TopicTime:time.Now()}

	//查看是否有重复的name字段
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}

	//插入字段
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}

	return nil
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		 return err
	}

	o := orm.NewOrm()
	cate := &Category{Id:cid}
	_, err = o.Delete(cate)

	return err
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}


func AddTopic(title, category, content string) error {
	o := orm.NewOrm()

	topic := &Topic{
		Title:title,
		Content:content,
		Category:category,
		Created:time.Now(),
		Updated:time.Now(),
		ReplyTime:time.Now(),
	}

	_, err := o.Insert(topic)

	return err
}

/*获取文章信息，isDesc : 是否为倒叙*/
func GetAllTopics(isDesc bool) (topics []*Topic,err error) {
	o := orm.NewOrm()

	topics = make([]*Topic , 0)

	qs := o.QueryTable("topic")

	if isDesc {
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}

	return topics, err
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()

	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)
	return topic, nil
}

func ModifyTopic(tid, title, category, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.Title = title
		topic.Content = content
		topic.Category = category
		topic.Updated = time.Now()
		o.Update(topic)
	}
	return nil
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	topic := &Topic{Id: tidNum}
	_, err = o.Delete(topic)
	return err
}