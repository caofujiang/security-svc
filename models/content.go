package models

import (
	"github.com/jinzhu/gorm"
	"github.com/secrity-svc/pkg/setting"
	gorm2 "gorm.io/gorm"
	"time"
)

type Content struct {
	Model
	Uid         string    `json:"uid"`
	Pid         int       `json:"pid"`
	Result      string    `json:"result"`
	IsDestroyed int       `json:"is_destroyed"`
	IsEnd       int       `json:"is_end"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
}

type Security_content struct {
	Model
	Uid         string    `json:"uid"`
	Pid         int       `json:"pid"`
	Result      string    `json:"result"`
	IsDestroyed int       `json:"is_destroyed"`
	IsEnd       int       `json:"is_end"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
}

// AddContent Add a Content
func AddContent(uid string, pid int, result string) error {

	if setting.DatabaseSetting.Type == "mysql" {
		content := Content{
			Uid:         uid,
			Pid:         pid,
			Result:      result,
			IsDestroyed: 0,
			IsEnd:       0,
			StartAt:     time.Now().Local(),
		}
		if err := db.Create(&content).Error; err != nil {
			return err
		}
	} else if setting.DatabaseSetting.Type == "cmdb" {
		securityContent := Security_content{
			Uid:         uid,
			Pid:         pid,
			Result:      result,
			IsDestroyed: 0,
			IsEnd:       0,
			StartAt:     time.Now().Local(),
		}
		if err := db2.Create(&securityContent).Error; err != nil {
			return err
		}
	}
	return nil
}

// ExistContentByUID determines whether a Content exists based on the UID
func ExistContentByUID(uid string) (Content, error) {
	var content Content
	var securityContent Security_content

	if setting.DatabaseSetting.Type == "mysql" {
		err := db.Select("result").Where("uid = ? ", uid).First(&content).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return content, err
		}

		return content, nil

	} else if setting.DatabaseSetting.Type == "cmdb" {
		err := db2.Select("result").Where("uid = ? ", uid).First(&securityContent).Error
		if err != nil && err != gorm2.ErrRecordNotFound {
			return Content(securityContent), err
		}
		return Content(securityContent), nil
	}

	return content, nil

}

// EditContent modify a single Content
func EditContent(uid string, data interface{}) error {
	if setting.DatabaseSetting.Type == "mysql" {
		if err := db.Model(&Content{}).Where("uid = ? ", uid).Updates(data).Error; err != nil {
			return err
		}
	} else if setting.DatabaseSetting.Type == "cmdb" {
		if err := db2.Model(&Security_content{}).Where("uid = ? ", uid).Updates(data).Error; err != nil {
			return err
		}
	}
	return nil
}
