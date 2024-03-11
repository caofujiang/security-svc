package content_service

import (
	"github.com/secrity-svc/models"
	"time"
)

type Content struct {
	Uid         string
	Pid         int
	Result      string
	StartAt     time.Time
	EndAt       time.Time
	IsDestroyed int
	IsEnd       int
}

func (t *Content) Add() error {
	return models.AddContent(t.Uid, t.Pid, t.Result)
}

func (t *Content) Edit(uid, newContent string) error {
	data := make(map[string]interface{})
	content, err := models.ExistContentByUID(uid)
	if err != nil {
		return err
	}
	data["result"] = content.Result + newContent
	return models.EditContent(uid, data)
}

func (t *Content) EditIsEndStatus(uid, newContent string) error {
	data := make(map[string]interface{})
	data["is_end"] = 1
	data["end_at"] = time.Now().Local()
	content, err := models.ExistContentByUID(uid)
	if err != nil {
		return err
	}
	data["result"] = content.Result + newContent
	return models.EditContent(uid, data)
}

func (t *Content) EditStatus(uid string, status int) error {
	data := make(map[string]interface{})
	data["is_destroyed"] = status
	return models.EditContent(uid, data)
}
