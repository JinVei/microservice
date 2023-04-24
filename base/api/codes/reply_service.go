package codes

import (
	"github.com/jinvei/microservice/base/framework/codes"
)

const (
	ErrReplySvcInternel = codes.Code(1000201)
	ErrReplyCommentPage = codes.Code(1000202)
	ErrReplyCommentItem = codes.Code(1000203)
)

func init() {
	codes.Register(map[codes.Code]string{
		ErrReplySvcInternel: "internal error",
		ErrReplyCommentPage: "get comment page error",
		ErrReplyCommentItem: "get comment item error",
	})
}
