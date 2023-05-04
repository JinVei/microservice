package entity

import "time"

type CommentContent struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:'主键/2023-04-13'" json:"id"`
	Content        []byte    `gorm:"column:content;type:varchar(512);comment:'评论内容/2023-04-13'" json:"content"`
	IP             string    `gorm:"column:ip;type:varchar(20);comment:'IP/2023-04-13'" json:"ip"`
	Platform       int8      `gorm:"column:platform;comment:'发布平台/2023-04-13'" json:"platform"`
	Device         string    `gorm:"column:device;type:varchar(20);comment:'发布设备/2023-04-13'" json:"device"`
	State          uint64    `gorm:"column:state;not null;comment:'状态/0启用/1删除'" json:"state"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:'数据库创建时间'" json:"created_at"`
	CreatedBy      uint64    `gorm:"column:create_by;not null;comment:'创建者'" json:"create_by"`
	CreateTime     uint64    `gorm:"column:create_time;not null;comment:'创建时间'" json:"create_time"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'数据库修改时间'" json:"updated_at"`
	LastModifyBy   uint64    `gorm:"column:last_modify_by;comment:'最后修改者'" json:"last_modify_by"`
	LastModifyTime uint64    `gorm:"column:last_modify_time;comment:'最后修改时间'" json:"last_modify_time"`
}

// Set table name
func (CommentContent) TableName() string {
	return "comment_content"
}
