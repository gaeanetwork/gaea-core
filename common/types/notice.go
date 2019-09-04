package types

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// NoticeLevel notice level
type NoticeLevel uint32

const (
	// Info notice info level
	Info NoticeLevel = iota + 1
	// Warn notice warn level
	Warn
)

// Check check whether value is valid
func (level NoticeLevel) Check() bool {
	if level > Warn {
		return false
	}

	return true
}

// NoticeType notice type
type NoticeType uint32

const (
	// AssetCertificate asset certificate
	AssetCertificate NoticeType = iota + 1
	// Financial receivable bill
	Financial
	// Channel channel module notice
	Channel
	// Tx transaction module notice
	Tx
	// User user module notice
	User
)

// Check check whether value is valid
func (t NoticeType) Check() bool {
	if t > User {
		return false
	}

	return true
}

// Notice notice in the EBaaS
type Notice struct {
	ID         string      `bson:"_id" json:"id,omitempty"`
	Title      string      `bson:"title" json:"title,omitempty"`
	Content    string      `bson:"content" json:"content,omitempty"`
	ToUser     string      `bson:"to_user" json:"to_user,omitempty"`
	UserName   string      `bson:"user_name" json:"user_name,omitempty"`
	Level      NoticeLevel `bson:"level" json:"level,omitempty"`
	Type       NoticeType  `bson:"type" json:"type,omitempty"`
	IsRead     bool        `bson:"read" json:"read,omitempty"`
	CreateTime int64       `bson:"create_time" json:"create_time,omitempty"`
	ReadTime   int64       `bson:"read_time" json:"read_time,omitempty"`
}

// NewNotice new a notice
func NewNotice(title, content, toUser string, level NoticeLevel, noticeType NoticeType) (*Notice, error) {
	if len(title) == 0 {
		return nil, errors.New("Not specified title")
	}

	if len(content) == 0 {
		return nil, errors.New("Not specified content")
	}

	if len(toUser) == 0 {
		return nil, errors.New("Not specified toUser")
	}

	return &Notice{
		ID:         uuid.New().String(),
		Title:      title,
		Content:    content,
		ToUser:     toUser,
		Level:      level,
		Type:       noticeType,
		CreateTime: time.Now().Unix(),
	}, nil
}
