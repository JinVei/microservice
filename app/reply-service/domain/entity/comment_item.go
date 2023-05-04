package entity

import (
	"time"
)

type CommentItem struct {
	ID             uint64    `gorm:"primary_key;column:id;comment:'主键/2023-04-13'"`
	Subject        uint64    `gorm:"column:subject;not null;comment:'评论区id'"`
	Parent         uint64    `gorm:"column:parent;not null;comment:'父评/0代表根评论/2023-04-13'"`
	Floor          uint64    `gorm:"column:floor;comment:'楼层/2023-04-13'"`
	UserID         uint64    `gorm:"column:userid;comment:'用户ID/2023-04-13'"`
	ReplyTo        uint64    `gorm:"column:replyto;comment:'回复用户ID/2023-04-13'"`
	Like           uint64    `gorm:"column:like;comment:'赞/2023-04-13'"`
	Dislike        uint64    `gorm:"column:dislike;comment:'踩/2023-04-13'"`
	ReplyCnt       uint64    `gorm:"column:reply_cnt;comment:'回复数/2023-04-13'"`
	State          uint64    `gorm:"column:state;not null;comment:'状态/0启用/1删除'"`
	Seq            uint64    `gorm:"column:state;not null;comment:'序列号, 每次更新行时+1'"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:'数据库创建时间'"`
	CreatedBy      uint64    `gorm:"column:create_by;not null;comment:'创建者'"`
	CreateTime     uint64    `gorm:"column:create_time;not null;comment:'创建时间'"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'数据库修改时间'"`
	LastModifyBy   uint64    `gorm:"column:last_modify_by;comment:'最后修改者'"`
	LastModifyTime uint64    `gorm:"column:last_modify_time;comment:'最后修改时间'"`
	//	ContentID      uint64    `gorm:"column:content_id;not null;comment:'评论内容id'"`
}

// define table name
func (CommentItem) TableName() string {
	return "comment_item"
}

// // define index
// func (CommentIndex) Indexes() []gorm.Index {
// 	return []gorm.Index{
// 		{
// 			Name: "subject_parent_floor_createdat",
// 			Fields: []gorm.Expr{
// 				gorm.Expr{Column: "subject"},
// 				gorm.Expr{Column: "parent"},
// 				gorm.Expr{Column: "floor"},
// 				gorm.Expr{Column: "created_at"},
// 			},
// 		},
// 	}
// }
