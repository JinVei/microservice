package entity

import "time"

type CommentSubject struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:'主键/2023-04-13'" json:"-"`
	ObjType        uint64    `gorm:"column:obj_type;not null;comment:'与评论区关联的系统的类型'" json:"objType"`
	ObjID          uint64    `gorm:"column:obj_id;not null;comment:'与评论区关联的系统的id'" json:"objID"`
	Like           uint64    `gorm:"column:like;comment:'赞/2023-04-13'" json:"like,omitempty"`
	Dislike        uint64    `gorm:"column:dislike;comment:'踩/2023-04-13'" json:"dislike,omitempty"`
	ReplyCnt       uint64    `gorm:"column:reply_cnt;comment:'评论数/2023-04-13'" json:"replyCnt,omitempty"`
	State          uint64    `gorm:"column:state;not null;comment:'状态/0启用/1删除'" json:"state"`
	Seq            uint64    `gorm:"column:seq;comment:'序列号, 每次更新行时+1'" json:"seq,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:'数据库创建时间'" json:"createdAt"`
	CreateBy       uint64    `gorm:"column:create_by;not null;comment:'创建者'" json:"createBy"`
	CreateTime     uint64    `gorm:"column:create_time;not null;comment:'创建时间'" json:"createTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'数据库修改时间'" json:"-"`
	LastModifyBy   uint64    `gorm:"column:last_modify_by;comment:'最后修改者'" json:"lastModifyBy,omitempty"`
	LastModifyTime uint64    `gorm:"column:last_modify_time;comment:'最后修改时间'" json:"lastModifyTime,omitempty"`
}
