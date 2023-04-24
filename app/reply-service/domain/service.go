package domain

import (
	"github.com/jinvei/microservice/base/api/proto/v1/app"
)

// type IReplyCommentService interface {
// 	ListCommentPage(ctx context.Context, in *app.ListCommentPageReq, opts ...grpc.CallOption) (*app.ListCommentPageResp, error)
// 	PutComment(ctx context.Context, in *app.PutCommentReq, opts ...grpc.CallOption) (*app.PutCommentResp, error)
// 	CreateSubject(ctx context.Context, in *app.CreateSubjectReq, opts ...grpc.CallOption) (*app.CreateSubjectResp, error)
// 	GetSubject(ctx context.Context, in *app.GetSubjectReq, opts ...grpc.CallOption) (*app.GetSubjectResp, error)
// }

type IReplyCommentService app.ReplyCommentServiceServer
