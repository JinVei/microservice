package entity

import "time"

type CommentContent struct {
	Id             uint64    `xorm:"'id' bigint(20) unsigned notnull comment('评论 Index ID/2023-04-13') pk"`
	Content        string    `xorm:"'content' varchar(512) default null comment('评论内容/2023-04-13')"`
	Ip             string    `xorm:"'ip' varchar(20) default null comment('IP/2023-04-13')"`
	Platform       int8      `xorm:"'platform' tinyint(8) default null comment('发布平台/2023-04-13')"`
	Device         string    `xorm:"'device' varchar(20) default null comment('发布设备/2023-04-13')"`
	State          uint64    `xorm:"'state' bigint(20) unsigned notnull comment('状态/0启用/1删除')"`
	CreatedAt      time.Time `xorm:"'created_at' timestamp notnull created comment('数据库创建时间')"`
	CreateBy       uint64    `xorm:"'create_by' bigint(20) unsigned notnull comment('创建者')"`
	CreateTime     int64     `xorm:"'create_time' bigint(20) notnull created comment('创建时间')"`
	UpdatedAt      time.Time `xorm:"'updated_at' timestamp notnull updated comment('数据库修改时间')"`
	LastModifyBy   uint64    `xorm:"'last_modify_by' bigint(20) unsigned default null comment('最后修改者')"`
	LastModifyTime int64     `xorm:"'last_modify_time' updated comment('最后修改时间')"`
}

// Set table name
func (CommentContent) TableName() string {
	return "comment_content"
}
