package entity

import "time"

type CommentSubject struct {
	Id             uint64    `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null;auto_increment;comment:'主键/2023-04-13'"`
	ObjType        uint64    `gorm:"column:obj_type;type:bigint(20) unsigned;not null;comment:'与评论区关联的系统的类型'"`
	ObjId          uint64    `gorm:"column:obj_id;type:bigint(20) unsigned;not null;comment:'与评论区关联的系统的id'"`
	Like           *uint64   `gorm:"column:like;type:bigint(20);comment:'赞/2023-04-13'"`
	Hate           *uint64   `gorm:"column:hate;type:bigint(20);comment:'踩/2023-04-13'"`
	Count          *uint64   `gorm:"column:count;type:bigint(20);comment:'评论数/2023-04-13'"`
	State          uint64    `gorm:"column:state;type:bigint(20) unsigned;not null;comment:'状态/0启用/1删除'"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:'数据库创建时间'"`
	CreateBy       uint64    `gorm:"column:create_by;type:bigint(20) unsigned;not null;comment:'创建者'"`
	CreateTime     uint64    `gorm:"column:create_time;type:bigint(20) unsigned;not null;comment:'创建时间'"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'数据库修改时间'"`
	LastModifyBy   *uint64   `gorm:"column:last_modify_by;type:bigint(20) unsigned;comment:'最后修改者'"`
	LastModifyTime *uint64   `gorm:"column:last_modify_time;type:bigint(20) unsigned;comment:'最后修改时间'"`
}
