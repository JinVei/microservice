package domain

import (
	"context"
	"errors"

	"github.com/jinvei/microservice/app/reply-service/domain/entity"
	"github.com/jinvei/microservice/base/api/proto/v1/dto"
)

var (
	DBRecordNotFound = errors.New("record not found")
)

type IReplyCommentRepository interface {
	//ListCommentsByFloor(ctx context.Context, subject, parent int64, page, numPerPage int) ([]*entity.CommentIndex, error)
	ListCommentsPageIds(ctx context.Context, subject, parent uint64, floor, numPerPage int) ([]uint64, error)
	ListCommetContents(ctx context.Context, ids []uint64) ([]entity.CommentContent, error)
	//	ListCommetItemsPage(ctx context.Context, page uint64) ([]entity.CommentItem, error)
	ListCommetItem(ctx context.Context, ids []uint64) ([]entity.CommentItem, error)
	GetCommentLastFloor(ctx context.Context, subject, parent uint64) (uint64, error)
	CreateComment(ctx context.Context, subject, parent, floor, userid, replyto uint64, cc entity.CommentContent) (entity.CommentItem, entity.CommentContent, error)
	GetSubject(ctx context.Context, id uint64) (entity.CommentSubject, error)
	BatchSubmitComments(ctx context.Context, comments []*dto.ReplyComment) ([]entity.CommentItem, []entity.CommentContent, error)
	BatchIncrCommentCount(ctx context.Context, comments []entity.CountableItem) ([]*entity.CommentItem, []*entity.CommentContent, error)
	BatchIncrSubjectCount(ctx context.Context, comments []entity.CountableItem) ([]entity.CommentSubject, error)
}
