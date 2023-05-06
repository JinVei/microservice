package repository

import (
	"context"
	"time"

	"github.com/jinvei/microservice/base/framework/log"

	"github.com/jinvei/microservice/app/reply-service/domain"
	"github.com/jinvei/microservice/app/reply-service/domain/entity"
	"xorm.io/xorm"
)

var flog = log.Default

type ReplyCommentRepository struct {
	xorm     *xorm.Engine
	systemid uint64
}

func NewReplyCommentRepository(xorm *xorm.Engine) domain.IReplyCommentRepository {
	return &ReplyCommentRepository{
		xorm: xorm,
	}
}

// type IReplyCommentRepository interface {
func (repo *ReplyCommentRepository) ListCommentsPageIds(ctx context.Context, subject, parent uint64, floor, numPerPage int) ([]uint64, error) {
	ci := make([]entity.CommentItem, 0, numPerPage)

	pageOffset := (floor / numPerPage) * numPerPage
	maxFloor := ((floor / numPerPage) + 1) * numPerPage

	err := repo.xorm.Context(ctx).Where("subject = ? AND parent = ? AND floor <= ? ", subject, parent, maxFloor).
		Limit(numPerPage, pageOffset).Select("id").Find(&ci)

	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, numPerPage)
	for _, e := range ci {
		ids = append(ids, e.Id)
	}
	return ids, nil
}

func (repo *ReplyCommentRepository) ListCommetContents(ctx context.Context, ids []uint64) ([]entity.CommentContent, error) {
	cc := make([]entity.CommentContent, 0, len(ids))
	err := repo.xorm.Context(ctx).In("id", ids).Find(&cc)

	return cc, err
}

// func (repo *ReplyCommentRepository) ListCommetItemsPage(ctx context.Context, page uint64) ([]entity.CommentItem, error) {
// 	return nil, nil
// }

func (repo *ReplyCommentRepository) ListCommetItem(ctx context.Context, ids []uint64) ([]entity.CommentItem, error) {
	ci := make([]entity.CommentItem, 0, len(ids))
	err := repo.xorm.Context(ctx).In("id", ids).Find(&ci)

	return ci, err
}

func (repo *ReplyCommentRepository) GetCommentLastFloor(ctx context.Context, subject, parent uint64) (uint64, error) {
	ci := entity.CommentItem{}

	ok, err := repo.xorm.Context(ctx).Where("subject = ? AND parent = ?", subject, parent).Desc("floor").Limit(1).Get(&ci)

	if !ok {
		return 0, domain.DBRecordNotFound
	}

	return ci.Floor, err
}

func (repo *ReplyCommentRepository) CreateComment(ctx context.Context, subject, parent, floor, userid, replyto uint64, cc entity.CommentContent) (entity.CommentItem, entity.CommentContent, error) {
	ci := entity.CommentItem{
		Subject:    subject,
		Parent:     parent,
		Floor:      floor,
		UserId:     userid,
		Replyto:    replyto,
		CreateBy:   repo.systemid,
		CreateTime: time.Now().Unix(),
		CreatedAt:  time.Now(),
	}

	n, err := repo.xorm.Context(ctx).Insert(&ci)
	if err != nil {
		return entity.CommentItem{}, entity.CommentContent{}, err
	}
	cc.Id = ci.Id

	flog.Debug("xorm insert n", "n", n)

	n, err = repo.xorm.Context(ctx).Insert(&cc)
	if err != nil {
		return entity.CommentItem{}, entity.CommentContent{}, err
	}

	return ci, cc, nil
}

func (repo *ReplyCommentRepository) GetSubject(ctx context.Context, id uint64) (entity.CommentSubject, error) {
	sbj := entity.CommentSubject{}
	exist, err := repo.xorm.Context(ctx).Where("id = ?", id).Get(&sbj)
	if err != nil {
		return entity.CommentSubject{}, err
	}
	if !exist {
		return entity.CommentSubject{}, domain.DBRecordNotFound
	}
	return sbj, nil
}

// }
