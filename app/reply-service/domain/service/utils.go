package service

import (
	"github.com/jinvei/microservice/app/reply-service/domain/entity"
	"github.com/jinvei/microservice/base/api/proto/v1/dto"
)

func newCommentItemFromEntity(v *entity.CommentItem) *dto.ReplyCommentItem {
	d := &dto.ReplyCommentItem{
		ID:       int64(v.Id),
		Subject:  int64(v.Subject),
		Parent:   int64(v.Parent),
		Floor:    int64(v.Floor),
		UserID:   int64(v.UserId),
		Replyto:  int64(v.Replyto),
		Like:     int64(v.Like),
		Dislike:  int64(v.Dislike),
		ReplyCnt: int64(v.ReplyCnt),
		State:    int64(v.State),
	}
	return d
}

func newCommentContentFromEntity(v *entity.CommentContent) *dto.ReplyCommentContent {
	d := &dto.ReplyCommentContent{
		ID:       int64(v.Id),
		Content:  []byte(v.Content),
		IP:       v.Ip,
		Platform: int64(v.Platform),
		Device:   v.Device,
		State:    int64(v.State),
	}
	return d
}

func newSubjectFromEntity(v *entity.CommentSubject) *dto.ReplyCommentSubject {
	d := &dto.ReplyCommentSubject{
		ID:       int64(v.Id),
		ObjType:  int64(v.ObjType),
		ObjID:    int64(v.ObjId),
		Like:     int64(v.Like),
		Dislike:  int64(v.Dislike),
		ReplyCnt: int64(v.ReplyCnt),
		Seq:      int64(v.Seq),
		State:    int64(v.State),
	}
	return d
}
