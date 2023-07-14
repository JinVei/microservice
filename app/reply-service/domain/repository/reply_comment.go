package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/jinvei/microservice/base/framework/log"
	"gorm.io/gorm"

	"github.com/jinvei/microservice/app/reply-service/domain"
	"github.com/jinvei/microservice/app/reply-service/domain/entity"
)

var flog = log.Default

type ReplyCommentRepository struct {
	orm      *gorm.DB
	systemid uint64
}

func NewReplyCommentRepository(orm *gorm.DB, systemid uint64) domain.IReplyCommentRepository {
	return &ReplyCommentRepository{
		orm:      orm,
		systemid: systemid,
	}
}

// type IReplyCommentRepository interface {
func (repo *ReplyCommentRepository) ListCommentsPageIds(ctx context.Context, subject, parent uint64, floor, numPerPage int) ([]uint64, error) {
	ci := make([]entity.CommentItem, 0, numPerPage)

	pageOffset := (floor / numPerPage) * numPerPage
	maxFloor := ((floor / numPerPage) + 1) * numPerPage

	dbc := repo.orm.WithContext(ctx).Where("subject = ? AND parent = ? AND floor <= ? ", subject, parent, maxFloor).
		Offset(pageOffset).Limit(numPerPage).Select("id").Find(&ci)

	if dbc.Error != nil {
		return nil, dbc.Error
	}
	ids := make([]uint64, 0, numPerPage)
	for _, e := range ci {
		ids = append(ids, e.Id)
	}
	return ids, nil
}

func (repo *ReplyCommentRepository) ListCommetContents(ctx context.Context, ids []uint64) ([]entity.CommentContent, error) {
	cc := make([]entity.CommentContent, 0, len(ids))
	res := repo.orm.WithContext(ctx).Where("id IN ?", ids).Find(&cc)

	return cc, res.Error
}

// func (repo *ReplyCommentRepository) ListCommetItemsPage(ctx context.Context, page uint64) ([]entity.CommentItem, error) {
// 	return nil, nil
// }

func (repo *ReplyCommentRepository) ListCommetItem(ctx context.Context, ids []uint64) ([]entity.CommentItem, error) {
	ci := make([]entity.CommentItem, 0, len(ids))
	res := repo.orm.WithContext(ctx).Where("id in ?", ids).Find(&ci)

	return ci, res.Error
}

func (repo *ReplyCommentRepository) GetCommentLastFloor(ctx context.Context, subject, parent uint64) (uint64, error) {
	ci := entity.CommentItem{}

	res := repo.orm.WithContext(ctx).Where("subject = ? AND parent = ?", subject, parent).Order("floor desc").Limit(1).Find(&ci)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) || ci.Id == 0 {
		// not record found, so this is the first floor
		return 0, nil
	}

	return ci.Floor, res.Error
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

	res := repo.orm.WithContext(ctx).Create(&ci)
	if res.Error != nil {
		return entity.CommentItem{}, entity.CommentContent{}, res.Error
	}
	cc.Id = ci.Id

	flog.Debug("CreateComment insert row", "row", res.RowsAffected)

	res = repo.orm.WithContext(ctx).Create(&cc)
	if res.Error != nil {
		return entity.CommentItem{}, entity.CommentContent{}, res.Error
	}

	return ci, cc, nil
}

func (repo *ReplyCommentRepository) GetSubject(ctx context.Context, id uint64) (entity.CommentSubject, error) {
	sbj := entity.CommentSubject{}
	res := repo.orm.WithContext(ctx).Where("id = ?", id).Find(&sbj)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return entity.CommentSubject{}, domain.DBRecordNotFound
	}

	if res.Error != nil {
		return entity.CommentSubject{}, res.Error
	}
	return sbj, nil
}

func (repo *ReplyCommentRepository) BatchSubmitComments(ctx context.Context, comments []*dto.ReplyComment) ([]entity.CommentItem, []entity.CommentContent, error) {
	items := make([]entity.CommentItem, 0, len(comments))
	contents := make([]entity.CommentContent, 0, len(comments))
	for _, i := range comments {
		ci := entity.CommentItem{
			Subject:    uint64(i.Item.Subject),
			Parent:     uint64(i.Item.Parent),
			Floor:      uint64(i.Item.Floor),
			UserId:     uint64(i.Item.UserID),
			Replyto:    uint64(i.Item.ReplyCnt),
			CreateBy:   repo.systemid,
			CreateTime: time.Now().Unix(),
			CreatedAt:  time.Now(),
		}

		contents = append(contents, entity.CommentContent{
			Content:    string(i.Content.Content),
			Ip:         i.Content.IP,
			Platform:   int8(i.Content.Platform),
			Device:     i.Content.Device,
			State:      uint64(i.Content.State),
			CreateBy:   repo.systemid,
			CreateTime: time.Now().Unix(),
			CreatedAt:  time.Now(),
		})
		items = append(items, ci)
	}

	res := repo.orm.WithContext(ctx).Create(&items)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	for i := 0; i < len(items); i++ {
		contents[i].Id = items[i].Id
	}

	res = repo.orm.WithContext(ctx).Create(&contents)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	return items, contents, nil
}

func (repo *ReplyCommentRepository) BatchIncrCommentCount(ctx context.Context, comments []entity.CountableItem) ([]*entity.CommentItem, []*entity.CommentContent, error) {
	groupBatchByLike := make(map[int][]uint64)
	groupBatchByReplycnt := make(map[int][]uint64)
	ids := make([]uint64, 0, len(comments))

	for _, v := range comments {
		groupBatchByLike[v.Like] = append(groupBatchByLike[v.Like], v.Id)
		groupBatchByReplycnt[v.Reply] = append(groupBatchByReplycnt[v.Reply], v.Id)
		ids = append(ids, v.Id)
	}

	for cnt, ids := range groupBatchByLike {
		res := repo.orm.Model(&entity.CommentItem{}).WithContext(ctx).Where(" id IN ?", ids).Updates(map[string]interface{}{"like_cnt": gorm.Expr("like_cnt + ?", cnt), "seq": gorm.Expr("seq + ?", 1), "last_modify_by": repo.systemid})
		if res.Error != nil {
			flog.Error(res.Error, `res := repo.orm.WithContext(ctx).Where(" id IN ?", ids).Update("like_cnt", gorm.Expr("like_cnt + ?", cnt)) error`, "ids", ids, "cnt", cnt)
		}
	}

	for cnt, ids := range groupBatchByReplycnt {
		res := repo.orm.Model(&entity.CommentItem{}).WithContext(ctx).Where(" id IN ?", ids).Updates(map[string]interface{}{"reply_cnt": gorm.Expr("reply_cnt + ?", cnt), "seq": gorm.Expr("seq + ?", 1), "last_modify_by": repo.systemid})
		if res.Error != nil {
			flog.Error(res.Error, `res := repo.orm.WithContext(ctx).Where(" id IN ?", ids).Update("reply_cnt", gorm.Expr("reply_cnt + ?", cnt)) error`, "ids", ids, "cnt", cnt)
		}
	}

	items := make([]*entity.CommentItem, 0, len(comments))
	contents := make([]*entity.CommentContent, 0, len(comments))

	res := repo.orm.Find(&items, ids)
	if res.Error != nil {
		flog.Error(res.Error, "repo.orm.Find(items, ids) error", "ids", ids)
		return nil, nil, res.Error
	}

	res = repo.orm.Find(&contents, ids)
	if res.Error != nil {
		flog.Error(res.Error, "repo.orm.Find(items, ids) error", "ids", ids)
		return nil, nil, res.Error
	}

	return items, contents, nil
}
func (repo *ReplyCommentRepository) BatchIncrSubjectCount(ctx context.Context, comments []entity.CountableItem) ([]entity.CommentSubject, error) {
	groupBatchByLike := make(map[int][]uint64)
	groupBatchByReplycnt := make(map[int][]uint64)
	ids := make([]uint64, 0, len(comments))

	for _, v := range comments {
		if 0 < v.Like {
			groupBatchByLike[v.Like] = append(groupBatchByLike[v.Like], v.Id)
		}
		if 0 < v.Reply {
			groupBatchByReplycnt[v.Reply] = append(groupBatchByReplycnt[v.Reply], v.Id)
		}
		ids = append(ids, v.Id)
	}

	for cnt, ids := range groupBatchByLike {
		res := repo.orm.Model(&entity.CommentSubject{}).WithContext(ctx).Where(" id IN ?", ids).Updates(map[string]interface{}{"like_cnt": gorm.Expr("like_cnt + ?", cnt), "seq": gorm.Expr("seq + ?", 1), "last_modify_by": repo.systemid})
		if res.Error != nil {
			flog.Error(res.Error, `res := repo.orm.WithContext(ctx).Where(" id IN ?", ids).Update("like_cnt", gorm.Expr("like_cnt + ?", cnt)) error`, "ids", ids, "cnt", cnt)
		}
	}

	for cnt, ids := range groupBatchByReplycnt {
		res := repo.orm.Model(&entity.CommentSubject{}).WithContext(ctx).Where(" id IN ?", ids).Updates(map[string]interface{}{"reply_cnt": gorm.Expr("reply_cnt + ?", cnt), "seq": gorm.Expr("seq + ?", 1), "last_modify_by": repo.systemid}) //Update("reply_cnt", gorm.Expr("reply_cnt + ?", cnt))
		if res.Error != nil {
			flog.Error(res.Error, `res := repo.orm.WithContext(ctx).Where(" id IN ?", ids).Update("reply_cnt", gorm.Expr("reply_cnt + ?", cnt)) error`, "ids", ids, "cnt", cnt)
		}
	}

	subjects := make([]entity.CommentSubject, 0, len(comments))

	res := repo.orm.Find(&subjects, ids)
	if res.Error != nil {
		flog.Error(res.Error, "repo.orm.Find(items, ids) error", "ids", ids)
		return nil, res.Error
	}

	return subjects, nil
}
