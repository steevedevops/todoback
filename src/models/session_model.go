package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Session struct {
	SessionID  int       `gorm:"primary_key"`
	SessionKey string    `json:"session_key"`
	UserAgent  string    `json:"user_agent"`
	ClientIp   string    `json:"client_ip"`
	IsBlocked  bool      `json:"is_blocked"`
	ExpiresAt  time.Time `json:"expires_at"`
	UserID     int
}

func (ss *Session) Save(db *gorm.DB) (*Session, error) {
	err := db.Debug().Create(&ss).Error

	if err != nil {
		return &Session{}, err
	}
	return ss, nil
}

func (ss *Session) UpdateSession(db *gorm.DB, columns map[string]interface{}, sessionId int) (*Session, error) {
	err := db.Debug().Model(&Session{}).Where("session_id = ?", sessionId).Take(&Session{}).UpdateColumns(
		columns,
	).Error
	if err != nil {
		return &Session{}, err
	}
	return ss, nil
}

func (ss *Session) IsExpired() bool {
	return ss.ExpiresAt.Before(time.Now())
}

func (ss *Session) FindSessionkey(db *gorm.DB, SessionID string, UserId int) (*Session, error) {
	err := db.Debug().Model(&Session{}).Where("session_key = ? AND user_id = ? AND is_blocked = ?", SessionID, UserId, false).Take(&ss).Error

	if err != nil {
		return &Session{}, err
	}
	return ss, nil
}

func (ss *Session) FindSessionUser(db *gorm.DB, UserId, SessionKey string) bool {
	fmt.Println(SessionKey)
	err := db.Debug().Model(&Session{}).Where("user_id = ? AND session_key = ? AND is_blocked = ?", UserId, SessionKey, false).Take(&ss).Error
	return err != nil
}

func (ss *Session) Delete(db *gorm.DB, SessionKey string) (int64, error) {
	db = db.Debug().Model(&Session{}).Where("session_key = ?", SessionKey).Take(&Session{}).Delete(&Session{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
