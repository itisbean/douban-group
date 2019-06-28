package model

import (
	"log"
	"strconv"

	"go-crawler/douban-group/parse"

	"github.com/jinzhu/gorm"
)

// Add 新增数据
func Add(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		if err := DB.Create(&topic).Error; err != nil {
			log.Printf("db.Create index: %s, err : %v", strconv.Itoa(index), err)
		}
	}
}

// GetVersion 获取上一次的页数
func GetVersion() int {
	var item parse.DoubanGroupDbhyz
	err := DB.Select("version").Order("version desc").First(&item).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return 0
	}
	return item.Version
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
