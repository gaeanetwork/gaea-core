package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Action int

const (
	// Login login ebaas
	Login Action = iota + 1
	// Logout logout ebaas
	Logout
	// SyncBlock sync blocks and transactions to MongoDB
	SyncBlock
	// PeerJoinChannel peer join channel
	PeerJoinChannel
)

var actions = [...]string{
	"Login",
	"Logout",
	"SyncBlock",
	"PeerJoinChannel",
}

// String returns the show name by the index
func (a Action) String() string {
	if Login <= a && a <= PeerJoinChannel {
		return actions[a-1]
	}
	return ""
}

type ActionLog struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	Action    Action `bson:"action" json:"action,omitempty"`
	TimeStamp int64  `bson:"timestamp" json:"timestamp,omitempty"`
	UserID    string `bson:"user_id" json:"user_id,omitempty"`
	SessionID string `bson:"session_id" json:"session_id,omitempty"`
	Mark      string `bson:"mark" json:"mark,omitempty"`
}

func NewActionLog(action Action, userID string) (*ActionLog, error) {
	strAction := action.String()
	if len(strAction) == 0 {
		return nil, fmt.Errorf("action value err:%d", action)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to randomly generated UUID, err:%s", err.Error())
	}

	return &ActionLog{
		ID:        id.String(),
		Action:    action,
		TimeStamp: time.Now().Unix(),
		UserID:    userID,
	}, nil
}
