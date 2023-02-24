package entity

import "time"

type Users struct {
	Id             uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Username       string    `gorm:"column:username"`
	Password       string    `gorm:"column:password"`
	Telnumber      string    `gorm:"column:telnumber"`
	Email          string    `gorm:"column:email"`
	Salt           string    `gorm:"column:salt"`
	Gender         uint8     `gorm:"column:gender"`
	Status         uint8     `gorm:"column:status"` // 0:normal; 1:block;
	CreatedAt      time.Time `gorm:"column:created_at;DEFAULT:CURRENT_TIMESTAMP"`
	CreateBy       uint64    `gorm:"column:create_by"`
	CreateTime     uint64    `gorm:"column:create_time"`
	UpdatedAt      time.Time `gorm:"column:updated_at;DEFAULT:CURRENT_TIMESTAMP;ON UPDATE:CURRENT_TIMESTAMP"`
	LastModifyBy   uint64    `gorm:"column:last_modify_by"`
	LastModifyTime uint64    `gorm:"column:last_modify_time"`
}

const (
	UserStatusNormal = 0
	UserStatusBlock  = 1
)
