package domain

import (
	"context"

	"github.com/jinvei/microservice/app/reply-service/domain/entity"
)

type IReplyCommentRepository interface {
	//ListCommentsByFloor(ctx context.Context, subject, parent int64, page, numPerPage int) ([]*entity.CommentIndex, error)
	ListCommentsPageIds(ctx context.Context, subject, parent int64, page, numPerPage int) ([]uint64, error)
	ListCommetContents(ctx context.Context, ids []uint64) ([]*entity.CommentContent, error)
	ListCommetIndex(ctx context.Context, ids []uint64) ([]*entity.CommentItem, error)
}
