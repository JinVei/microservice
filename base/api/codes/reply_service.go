package codes

import (
	"github.com/jinvei/microservice/base/framework/codes"
)

const (
	ReplySvcSystemID = 10002

	ErrReplySvcInternel      = codes.Code(1000201)
	ErrReplySvcCommentPage   = codes.Code(1000202)
	ErrReplySvcCommentItem   = codes.Code(1000203)
	ErrReplySvcCreateComment = codes.Code(1000204)
	ErrReplySvcPutComment    = codes.Code(1000205)
	ErrReplySvcGetSubject    = codes.Code(1000206)
)

func init() {
	codes.Register(map[codes.Code]string{
		ErrReplySvcInternel:      "internal error",
		ErrReplySvcCommentPage:   "Get comment page error",
		ErrReplySvcCommentItem:   "Get comment item error",
		ErrReplySvcCreateComment: "Create Comment Error",
		ErrReplySvcPutComment:    "Put Comment Error",
		ErrReplySvcGetSubject:    "Get Subject Error",
	})
}
