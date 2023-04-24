package data

import (
	"trustPlatform/constant"
)

type SharedSM2Message struct {
	// couchDB使用的type
	ObjectType string `json:"docType"`

	Uid string `json:"uid"`
	// 加密内容
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	FileName  string `json:"fileName"`
	ToName    string `json:"toName"`
}

func NewSharedSM2Message(uid, content, timestamp, fileName, toName string) *SharedSM2Message {
	return &SharedSM2Message{
		ObjectType: constant.SharedSM2Message,
		Uid:        uid,
		Content:    content,
		Timestamp:  timestamp,
		FileName:   fileName,
		ToName:     toName,
	}
}
