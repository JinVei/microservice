package entity

import "time"

// type CommentSubject struct {
// 	Id             uint64    `xorm:"'id' pk autoincr comment('主键/2023-04-13') bigint(20) unsigned"`
// 	ObjType        uint64    `xorm:"'obj_type' notnull comment('与评论区关联的系统的类型') bigint(20) unsigned"`
// 	ObjId          uint64    `xorm:"'obj_id' notnull comment('与评论区关联的系统的id') bigint(20) unsigned"`
// 	Like           uint64    `xorm:"'like' null comment('赞/2023-04-13') bigint(20)"`
// 	Dislike        uint64    `xorm:"'dislike' null comment('踩/2023-04-13') bigint(20)"`
// 	ReplyCnt       uint64    `xorm:"'reply_cnt' null comment('评论数/2023-04-13') bigint(20)"`
// 	State          int64     `xorm:"'state' notnull comment('状态/0启用/1删除') bigint(20) unsigned"`
// 	Seq            uint64    `xorm:"'seq' null comment('序列号, 每次更新行时+1') bigint(20) unsigned"`
// 	CreatedAt      time.Time `xorm:"'created_at' notnull created comment('数据库创建时间')"`
// 	CreateBy       uint64    `xorm:"'create_by' notnull comment('创建者') bigint(20) unsigned"`
// 	CreateTime     uint64    `xorm:"'create_time' notnull created comment('创建时间') bigint(20) unsigned"`
// 	UpdatedAt      time.Time `xorm:"'updated_at' notnull updated comment('数据库修改时间')"`
// 	LastModifyBy   uint64    `xorm:"'last_modify_by' null comment('最后修改者') bigint(20) unsigned"`
// 	LastModifyTime uint64    `xorm:"'last_modify_time' null updated comment('最后修改时间') bigint(20) unsigned"`
// }

type CommentSubject struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement;column:id;comment:'主键/2023-04-13'"`
	ObjType        uint64    `gorm:"column:obj_type;not null;comment:'与评论区关联的系统的类型'"`
	ObjId          uint64    `gorm:"column:obj_id;not null;comment:'与评论区关联的系统的id'"`
	Like           uint64    `gorm:"column:like_cnt;comment:'赞/2023-04-13'"`
	Dislike        uint64    `gorm:"column:dislike;comment:'踩/2023-04-13'"`
	ReplyCnt       uint64    `gorm:"column:reply_cnt;comment:'评论数/2023-04-13'"`
	State          int64     `gorm:"column:state;not null;comment:'状态/0启用/1删除'"`
	Seq            uint64    `gorm:"column:seq;comment:'序列号, 每次更新行时+1'"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;autoCreateTime;comment:'数据库创建时间'"`
	CreateBy       uint64    `gorm:"column:create_by;not null;comment:'创建者'"`
	CreateTime     uint64    `gorm:"column:create_time;not null;autoCreateTime:nano;comment:'创建时间'"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;autoUpdateTime;comment:'数据库修改时间'"`
	LastModifyBy   uint64    `gorm:"column:last_modify_by;comment:'最后修改者'"`
	LastModifyTime uint64    `gorm:"column:last_modify_time;comment:'最后修改时间';autoUpdateTime:nano"`
}
