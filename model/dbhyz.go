package model

import (
	//"log"

	"go-crawler/douban-group/parse"

	"github.com/jinzhu/gorm"

	log "github.com/cihub/seelog"
)

// Add 新增数据
func Add(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		if err := DB.Create(&topic).Error; err != nil {
			log.Critical("create index: %d, err : %v", index, err)
		}
	}
}

// GetVersion 获取上一次的页数 （弃）
func GetVersion(min int, max int) (v []int) {
	var items []parse.DoubanGroupDbhyz
	err := DB.Select("version").Order("version").Group("version").Find(&items).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		log.Critical("get version err : %v", err)
		return
	}

	for i:=(min+1);i<=max;i++ {
		flag := true
		for _, item := range items {
			if i == item.Version {
				flag = false
			}
		}
		if flag == true {
			prev := i-1
			next := i+1
			if len(v) == 0 || i != v[len(v)-1] {
				if len(v) == 0 || (prev >0 && prev != v[len(v)-1]) {
					v = append(v, prev)
				}
				v = append(v, i)
			}
			if next <= max {
				v = append(v, next)
			}
		}
	}
	
	return v
}

// Save 新增或更新
func Save(topics []parse.DoubanGroupDbhyz) {
	for index, topic := range topics {
		err := DB.Where(parse.DoubanGroupDbhyz{TopicID: topic.TopicID}).Assign(topic).FirstOrCreate(&topic).Error
		if err != nil {
			log.Critical("first or create index: %d, err : %v", index, err)
		}
	}
}

// Update 更新数据
func Update(topic parse.DoubanGroupDbhyz) {
	err := DB.Model(&topic).Updates(topic).Error
	if err != nil {
		log.Critical("save err : %v", err)
	}
}

// GetNoContent 获取未更新详情页的贴子
func GetNoContent() (topics []parse.DoubanGroupDbhyz) {
	err := DB.Where("content IS NULL AND new_reply_time > '2019-01-01 00:00:00'").Find(&topics).Error
	if err != nil {
		log.Critical("get no content err : %v", err)
	}
	return
}

// GetOne 获取单条贴子数据
func GetOne(topicID int) (topic parse.DoubanGroupDbhyz) {
	err := DB.Where("topic_id = ?", topicID).First(&topic).Error
	if err != nil {
		log.Critical("get one err : %v", err)
	}
	return
}
