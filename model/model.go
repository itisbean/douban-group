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

	host = "127.0.0.1"
	port = "3306"
)

func init() {
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName))
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

// GetVersion 获取上一次的页数
func GetVersion() (int, error) {
	var item parse.DoubanGroupDbhyz
	err := DB.Select("version").Order("version desc").First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return item.Version, nil
}

// Save 新增或更新
func Save(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		err := DB.Where(parse.DoubanGroupDbhyz{TopicID: topic.TopicID}).Assign(topic).FirstOrCreate(&topic).Error
		if err != nil {
			log.Printf("index: %s, err : %v", strconv.Itoa(index), err)
		}
	}
}
