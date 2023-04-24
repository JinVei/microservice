package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/jinvei/microservice/app/reply-service/domain"
	"github.com/jinvei/microservice/app/reply-service/domain/entity"
	apicode "github.com/jinvei/microservice/base/api/codes"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/jinvei/microservice/base/framework/codes"
	"github.com/jinvei/microservice/base/framework/log"
	"golang.org/x/sync/singleflight"
)

var flog = log.Default

var _ domain.IReplyCommentService = &ReplyCommentService{}

const (
	pageGroupFromat = "subject:%d:parent:%d:page:%d"
)

type ReplyCommentService struct {
	repo      domain.IReplyCommentRepository
	cache     CommentCache
	pageGroup singleflight.Group
}

func (s *ReplyCommentService) ListCommentPage(ctx context.Context, in *app.ListCommentPageReq) (*app.ListCommentPageResp, error) {
	page := pageFromFloor(in.Floor)
	pageResourceKey := fmt.Sprintf(pageGroupFromat, in.Subject, in.Parent, page)

	// group page resource.
	// only one request would be executed for per page resouce where in case of concurrency
	resouce, err, shared := s.pageGroup.Do(pageResourceKey, func() (interface{}, error) {
		ids, err := s.getCommentPageIds(ctx, in.Subject, in.Parent, page)
		if err != nil {
			flog.Error(err, "getCommentPageIds(ctx, in.Subject, in.Parent, page)", "in.Subject", in.Subject, "in.Parent", in.Parent, "page", page)
			return apicode.ErrReplyCommentPage, nil
		}

		commentsM := make(map[uint64]*dto.ReplyComment, len(ids))

		if err := s.loadCommentPage(ctx, uint64(in.Subject), uint64(in.Parent), uint64(page), ids, commentsM); err != nil {
			flog.Errorf(err, "s.loadCommentPage(ctx,%v,%v,%v,%v,%v:%v)", uint64(in.Subject), uint64(in.Parent), uint64(page), ids, commentsM)
			return apicode.ErrReplyCommentItem, nil
		}

		comments, err := s.sortCommentToSlice(commentsM)
		if err != nil {
			flog.Error(err, "sortCommentToSlice()", "commentsM", commentsM)
			return apicode.ErrReplySvcInternel, nil
		}

		return comments, nil
	})

	if err != nil {
		return nil, err
	}

	errCode, ok := resouce.(codes.Code)
	if ok {
		return &app.ListCommentPageResp{
			Status: errCode.ToStatus(),
		}, nil
	}

	comments := resouce.([]*dto.ReplyComment)
	if shared {
		m := make([]*dto.ReplyComment, 0, len(comments))
		for _, v := range comments {
			m = append(m, proto.Clone(v).(*dto.ReplyComment))
		}
		comments = m
	}

	return &app.ListCommentPageResp{
		Status:   apicode.StatusOK.ToStatus(),
		Comments: comments,
	}, nil

}

func (s *ReplyCommentService) PutComment(ctx context.Context, in *app.PutCommentReq) (*app.PutCommentResp, error) {
	return nil, nil
}

func (s *ReplyCommentService) CreateSubject(ctx context.Context, in *app.CreateSubjectReq) (*app.CreateSubjectResp, error) {
	return nil, nil
}

func (s *ReplyCommentService) GetSubject(ctx context.Context, in *app.GetSubjectReq) (*app.GetSubjectResp, error) {
	return nil, nil
}

func (s *ReplyCommentService) sortCommentToSlice(m map[uint64]*dto.ReplyComment) ([]*dto.ReplyComment, error) {
	comments := make([]*dto.ReplyComment, 0, len(m))

	for _, v := range m {
		comments = append(comments, v)
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Item.Floor < comments[j].Item.Floor
	})

	return comments, nil
}

func (s *ReplyCommentService) loadCommentPage(ctx context.Context, subject, parent, page uint64, ids []uint64, comments map[uint64]*dto.ReplyComment) error {

	// load pageItem
	pageItem, err := s.cache.GetCommentIndexPage(ctx, subject, parent, page)
	if err == ErrCacheMiss {
		// load from db
		entitys, err := s.repo.ListCommetIndex(ctx, ids)
		if err != nil {
			return err
		}
		pageItem = make([]*dto.ReplyCommentItem, 0, len(entitys))
		for _, e := range entitys {
			it := newCommentIndexFromEntity(e)
			comments[uint64(it.ID)].Item = it
			pageItem = append(pageItem, it)
		}

		if err := s.cache.SoreCommentIndexPage(ctx, subject, parent, page, pageItem); err != nil {
			flog.Error(err, "s.cache.SoreCommentIndexPage(ctx, subject, parent, page, pageItem)",
				"subject", subject, "parent", parent, "page", page, "pageItem", pageItem)
		}
	} else {
		for _, it := range pageItem {
			comments[uint64(it.ID)].Item = it
		}
	}

	// TODO: merger comment like/count field to pageItem

	// load pageContent
	pageContent, err := s.cache.GetCommentContentPage(ctx, subject, parent, page)
	if err == ErrCacheMiss {
		// load from db
		entitys, err := s.repo.ListCommetContents(ctx, ids)
		if err != nil {
			return err
		}
		pageContent = make([]*dto.ReplyCommentContent, 0, len(entitys))
		for _, e := range entitys {
			it := newCommentContentFromEntity(e)
			comments[uint64(it.ID)].Content = it
			pageContent = append(pageContent, it)
		}

		if err := s.cache.SoreCommentContentPage(ctx, subject, parent, page, pageContent); err != nil {
			flog.Error(err, "s.cache.SoreCommentContentPage(ctx, subject, parent, page, pageContent)",
				"subject", subject, "parent", parent, "page", page, "pageItem", pageContent)
		}
	} else {
		for _, it := range pageContent {
			comments[uint64(it.ID)].Content = it
		}
	}

	return nil
}

// func (s *ReplyCommentService) loadCommentContentToComments(ctx context.Context, ids []uint64, comments map[uint64]*dto.ReplyComment) error {
// 	missIds := []uint64{}
// 	for i, idx := range ids {
// 		if _, exist := comments[idx]; !exist {
// 			comments[idx] = &dto.ReplyComment{}
// 		}

// 		content, err := s.cache.GetCommentContent(ctx, int64(idx))
// 		if err == ErrCacheMiss {
// 			missIds[i] = idx
// 			continue
// 		} else if err != nil {
// 			flog.Error(err, "cache.GetCommentContent() err:")
// 			continue
// 		}

// 		comments[idx].Content = content
// 	}

// 	if 0 < len(missIds) {
// 		idx, err := s.repo.ListCommetContents(ctx, missIds)
// 		if err != nil {
// 			flog.Error(err, "loadCommentContents() err")
// 		}

// 		for _, v := range idx {
// 			d := &dto.ReplyCommentContent{
// 				ID:       int64(v.ID),
// 				Content:  []byte(v.Content),
// 				IP:       v.IP,
// 				Platform: int64(v.Platform),
// 				Device:   v.Device,
// 				State:    int64(v.State),
// 			}
// 			comments[v.ID].Content = d
// 			if err := s.cache.StoreCommentContent(ctx, d); err != nil {
// 				flog.Error(err, "StoreCommentIndex()", "index", d)
// 			}
// 		}
// 	}

// 	return nil
// }

// func (s *ReplyCommentService) loadCommentIndexToComments(ctx context.Context, ids []uint64, comments map[uint64]*dto.ReplyComment) error {
// 	missIds := []uint64{}
// 	for i, idx := range ids {
// 		if _, exist := comments[idx]; !exist {
// 			comments[idx] = &dto.ReplyComment{}
// 		}

// 		index, err := s.cache.GetCommentIndex(ctx, idx)
// 		if err == ErrCacheMiss {
// 			missIds[i] = idx
// 			continue
// 		} else if err != nil {
// 			flog.Error(err, "cache.GetCommentContent() err:")
// 			continue
// 		}

// 		comments[idx].Index = index
// 	}

// 	// load miss's comment
// 	if 0 < len(missIds) {
// 		idx, err := s.repo.ListCommetIndex(ctx, missIds)
// 		if err != nil {
// 			flog.Error(err, "loadCommentContents() err")
// 		}

// 		for _, v := range idx {
// 			d := &dto.ReplyCommentItem{
// 				ID:      int64(v.ID),
// 				Subject: int64(v.Subject),
// 				Parent:  int64(v.Parent),
// 				Floor:   int64(v.Floor),
// 				UserID:  int64(v.UserID),
// 				Replyto: int64(v.ReplyTo),
// 				Like:    int64(v.Like),
// 				Hate:    int64(v.Hate),
// 				Count:   int64(v.Count),
// 				State:   int64(v.State),
// 			}
// 			comments[v.ID].Index = d
// 			if err := s.cache.StoreCommentIndex(ctx, d); err != nil {
// 				flog.Error(err, "StoreCommentIndex()", "index", d)
// 			}
// 		}
// 	}

// 	return nil
// }

const numPerPage = 5

// 5 comments per page. group by floor
func pageFromFloor(floor int64) int64 {
	page := floor / int64(numPerPage)
	return page
}

func (s *ReplyCommentService) getCommentPageIds(ctx context.Context, subject, parent, page int64) ([]uint64, error) {
	// get page from cache
	ids, err := s.cache.GetCommentPageIds(ctx, subject, parent, page)
	if err == ErrCacheMiss {
		// load index data from db to cache
		// if page hasn't data(cause by delete), skip next page.
		entitys, err := s.repo.ListCommentsPageIds(ctx, subject, parent, int(page), numPerPage)
		if err != nil {
			return nil, err
		}
		ids = entitys
		if err := s.cache.StoreCommentPageIds(ctx, subject, parent, page, ids); err != nil {
			flog.Error(err, "s.cache.StoreCommentPageIds()", "subject", "parent", parent, "page", page, "ids", ids)
		}
	} else if err != nil {
		// TODO check end
		return nil, err
	}

	return ids, nil
}

func newCommentIndexFromEntity(v *entity.CommentItem) *dto.ReplyCommentItem {
	d := &dto.ReplyCommentItem{
		ID:      int64(v.ID),
		Subject: int64(v.Subject),
		Parent:  int64(v.Parent),
		Floor:   int64(v.Floor),
		UserID:  int64(v.UserID),
		Replyto: int64(v.ReplyTo),
		Like:    int64(v.Like),
		Hate:    int64(v.Dislike),
		Count:   int64(v.Count),
		State:   int64(v.State),
	}
	return d
}

func newCommentContentFromEntity(v *entity.CommentContent) *dto.ReplyCommentContent {
	d := &dto.ReplyCommentContent{
		ID:       int64(v.ID),
		Content:  v.Content,
		IP:       v.IP,
		Platform: int64(v.Platform),
		Device:   v.Device,
		State:    int64(v.State),
	}
	return d
}
