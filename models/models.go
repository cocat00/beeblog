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

// 评论
type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB()  {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册定义的 model
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
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
	if err != nil {
		return err
	}

	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		// 如果不存在我们就直接忽略，只当分类存在时进行更新
		cate.TopicCount++
		_, err = o.Update(cate)
		if err != nil {
			return err
		}
	}

	return err
}

/*获取文章信息，isDesc : 是否为倒叙*/
func GetAllTopics(category string, isDesc bool) (topics []*Topic,err error) {
	o := orm.NewOrm()

	topics = make([]*Topic , 0)

	qs := o.QueryTable("topic")

	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}

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

	//临时变量
	var oldCate string

	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		topic.Title = title
		topic.Content = content
		topic.Category = category
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}

	//更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
			if err != nil {
				return err
			}
		}
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	//更新分类统计
	var oldCate string
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", topic.Category).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
			if err != nil {
				return err
			}
		}
	}

	return err
}


func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	return err
}

func GetAllReplies(tid string) (replies []*Comment, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Comment, 0)

	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	reply := &Comment{Id:ridNum}

	_, err = o.Delete(reply)
	return err
}