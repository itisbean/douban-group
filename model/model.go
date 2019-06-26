package model

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql

	"go-crawler/douban-group/parse"
)

// mysql config
var (
	DB *gorm.DB

	username = "root"
	password = "dony123."
	dbName   = "spiders"

	host = "192.168.254.128"
	port = "3306"
)

func init() {
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName))
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}

	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "sp_" + defaultTableName
	}
}

// Add 新增数据
func Add(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		if err := DB.Create(&topic).Error; err != nil {
			log.Printf("db.Create index: %s, err : %v", strconv.Itoa(index), err)
		}
	}
}

// GetLastTime 获取最新时间
func GetLastTime() (string, error) {
	time := ""
	err := DB.Select("new_reply_time").Order("new_reply_time desc").First(&time).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return time, err
	}
	return time, nil
}

// Save 新增或更新
func Save(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		err := DB.Where(parse.DoubanGroupDbhyz{TopicID: topic.TopicID}).Assign(topic).FirstOrCreate(&topic).Error
		if err != nil {
			log.Printf("db.Create index: %s, err : %v", strconv.Itoa(index), err)
		}
	}
}
