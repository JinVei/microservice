package entity

import (
	"time"
)

// type CommentItem struct {
// 	Id             uint64    `xorm:"'id' pk autoincr unsigned bigint(20)"`    // 主键/2023-04-13
// 	Subject        uint64    `xorm:"'subject' notnull unsigned bigint(20)"`   // 评论区id
// 	Parent         uint64    `xorm:"'parent' notnull unsigned bigint(20)"`    // 父评/0代表根评论/2023-04-13
// 	Floor          uint64    `xorm:"'floor' bigint(20)"`                      // 楼层/2023-04-13
// 	UserId         uint64    `xorm:"'user_id' bigint(20)"`                    // 用户ID/2023-04-13
// 	Replyto        uint64    `xorm:"'replyto' bigint(20)"`                    // 回复用户ID/2023-04-13
// 	Like           uint64    `xorm:"'like' bigint(20)"`                       // 赞/2023-04-13
// 	Dislike        uint64    `xorm:"'dislike' bigint(20)"`                    // 踩/2023-04-13
// 	ReplyCnt       uint64    `xorm:"'reply_cnt' bigint(20)"`                  // 回复数/2023-04-13
// 	State          uint64    `xorm:"'state' notnull unsigned bigint(20)"`     // 状态/0启用/1删除
// 	Seq            uint64    `xorm:"'seq' bigint(20)"`                        // 序列号, 每次更新行时+1
// 	CreatedAt      time.Time `xorm:"'created_at' created"`                    // 数据库创建时间
// 	CreateBy       uint64    `xorm:"'create_by' notnull unsigned bigint(20)"` // 创建者
// 	CreateTime     int64     `xorm:"'create_time' created"`                   // 创建时间
// 	UpdatedAt      time.Time `xorm:"'updated_at' updated"`                    // 数据库修改时间
// 	LastModifyBy   uint64    `xorm:"'last_modify_by' bigint(20)"`             // 最后修改者
// 	LastModifyTime uint64    `xorm:"'last_modify_time' updated bigint(20)"`   // 最后修改时间
// }

type CommentItem struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Subject        uint64    `gorm:"column:subject"`
	Parent         uint64    `gorm:"column:parent"`
	Floor          uint64    `gorm:"column:floor"`
	UserId         uint64    `gorm:"column:user_id"`
	Replyto        uint64    `gorm:"column:replyto"`
	Like           uint64    `gorm:"column:like_cnt"`
	Dislike        uint64    `gorm:"column:dislike"`
	ReplyCnt       uint64    `gorm:"column:reply_cnt"`
	State          uint64    `gorm:"column:state"`
	Seq            uint64    `gorm:"column:seq"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	CreateBy       uint64    `gorm:"column:create_by"`
	CreateTime     int64     `gorm:"column:create_time;autoCreateTime:nano"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	LastModifyBy   uint64    `gorm:"column:last_modify_by"`
	LastModifyTime uint64    `gorm:"column:last_modify_time;autoUpdateTime:nano"`
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
