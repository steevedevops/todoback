package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/steevepypo/todoback/src/services/security"
)

type User struct {
	UserID     int       `json:"user_id" gorm:"primary_key;auto_increment"`
	Username   string    `json:"username" gorm:"unique;size:150;not null"`
	FirstName  string    `json:"first_name" gorm:"size:150"`
	LastName   string    `json:"last_name" gorm:"size:150"`
	Email      string    `json:"email" gorm:"size:255"`
	Password   string    `json:"password" gorm:"size:100;not null;"`
	DateJoined time.Time `json:"date_joined" gorm:"default:CURRENT_TIMESTAMP"`
	LastLogin  time.Time `json:"last_login" gorm:"default:CURRENT_TIMESTAMP"`
	IsActive   bool      `json:"is_active"`
	// Para funcinoar um foreigKey tem que existir um primary key na outra tabela
	Sessions []Session `json:"sessions" gorm: "ForeignKey:UserID"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.LastLogin = time.Now()
	// u.UpdatedAt = time.Now()
}

func (u *User) Save(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) List(db *gorm.DB, UserId int) (*[]User, error) {
	users := []User{}
	err := db.Debug().Model(&User{}).Where("user_id = ?", UserId).Limit(10000).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

func (u *User) FindById(db *gorm.DB, UserId string) (*User, error) {
	err := db.Debug().Model(&User{}).Where("user_id = ?", UserId).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) Delete(db *gorm.DB) error {
	users := []User{}
	err := db.Debug().Model(&User{}).Delete(&users).Error
	if err != nil {
		return err
	}
	return nil
}
