package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jinvei/microservice/app/reply-service/domain"
	apicode "github.com/jinvei/microservice/base/api/codes"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/jinvei/microservice/base/framework/cache"
	"github.com/jinvei/microservice/base/framework/codes"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/log"
	"golang.org/x/sync/singleflight"
)

var flog = log.Default

var _ domain.IReplyCommentService = &ReplyCommentService{}

const (
	pageGroupFromat      = "comments:subject:%d:parent:%d:page:%d"
	lastFloorGroupFromat = "last-floor:subject:%d:parent:%d"
	subjectGroupFromat   = "subject:%d"
)

type ReplyCommentService struct {
	repo         domain.IReplyCommentRepository
	cache        *CommentCache
	resouceGroup singleflight.Group
	systemid     int64

	batchWriter *CommentBatchWriter
}

func NewReplyCommentService(repo domain.IReplyCommentRepository, conf configuration.Configuration) *ReplyCommentService {
	systemid := conf.GetSystemID()
	sid, _ := strconv.Atoi(systemid)

	if systemid == "" {
		panic("systemid is empty. should setting SystemID by configuration.SetSystemID(")
	}
	rbd := cache.RedisClient(conf)

	svcCfg := defaultReplyCommentSvcConfig()
	err := conf.GetSvcJson(systemid, "", &svcCfg)
	if err != nil {
		flog.Warn("NewReplyCommentService:conf.GetSvcJson() error", "err", err)
	}

	flog.Info("ReplyCommentService Cfg", "cfg", svcCfg)

	cache := NewCommentCache(rbd, svcCfg.CacheDura, svcCfg.IndexCacheDura)

	bd, err := time.ParseDuration(svcCfg.BatchWriteInterval)
	if err != nil {
		flog.Error(err, "time.ParseDuration(svcCfg.BatchWriteInterval) err", "svcCfg.BatchWriteInterval", svcCfg.BatchWriteInterval)
	}
	bd = 5 * time.Second

	batchWriter := NewCommentBatchWriter(repo, cache, svcCfg.BatchWriteNum, bd)
	batchWriter.Start()

	return &ReplyCommentService{
		repo:     repo,
		cache:    cache,
		systemid: int64(sid),

		batchWriter: batchWriter,
	}
}

func (s *ReplyCommentService) ListCommentPage(ctx context.Context, in *app.ListCommentPageReq) (*app.ListCommentPageResp, error) {
	flog.Info("ListCommentPage", "req", in)

	page := pageFromFloor(in.Floor)
	pageResourceKey := fmt.Sprintf(pageGroupFromat, in.Subject, in.Parent, page)

	// group page resource.
	// only one request would be executed for per page resouce where in case of concurrency
	resouce, err, shared := s.resouceGroup.Do(pageResourceKey, func() (interface{}, error) {
		ids, err := s.getCommentPageIds(ctx, uint64(in.Subject), uint64(in.Parent), uint64(page), uint64(in.Floor))
		if err != nil {
			flog.Error(err, "getCommentPageIds(ctx, in.Subject, in.Parent, page)", "in.Subject", in.Subject, "in.Parent", in.Parent, "page", page)
			return apicode.ErrReplySvcCommentPage, nil
		}

		commentsM := make(map[uint64]*dto.ReplyComment)

		if err := s.loadCommentItemToComments(ctx, ids, commentsM); err != nil {
			flog.Errorf(err, "s.loadCommentIndexToComments(ctx,%v,%v):", ids, commentsM)
			return apicode.ErrReplySvcCommentItem, nil
		}

		if err := s.loadCommentContentToComments(ctx, ids, commentsM); err != nil {
			flog.Errorf(err, "s.loadCommentContentToComments(ctx,%v,%v):", ids, commentsM)
			return apicode.ErrReplySvcCommentItem, nil
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
	// get last floor
	floor, err := s.getSubjectLastFloor(ctx, uint64(in.Subject), uint64(in.Parent))
	if err != nil {
		return &app.PutCommentResp{Status: apicode.ErrReplySvcInternel.ToStatus()}, nil
	}

	// insert comment to db
	err = s.batchWriter.PutComment(&dto.ReplyComment{
		Item: &dto.ReplyCommentItem{
			Subject: in.Subject,
			Parent:  in.Parent,
			Floor:   floor,
			UserID:  in.UserID,
			Replyto: in.ReplyTo,
		},
		Content: in.Content,
	})
	if err != nil {
		flog.Error(err, "PutComment()", "in.Subject", in.Subject, "in.Parent", in.Parent, "floor", floor)
		return &app.PutCommentResp{
			Status: apicode.ErrReplySvcCreateComment.ToStatus(),
		}, nil
	}

	return &app.PutCommentResp{Status: apicode.StatusOK.ToStatus()}, nil
}

// TODO
func (s *ReplyCommentService) CreateSubject(ctx context.Context, in *app.CreateSubjectReq) (*app.CreateSubjectResp, error) {
	return nil, nil
}

// TODO
func (s *ReplyCommentService) GetSubject(ctx context.Context, in *app.GetSubjectReq) (*app.GetSubjectResp, error) {
	// Get From Cache
	resouceKey := fmt.Sprintf(subjectGroupFromat, in.ID)
	res, err, shared := s.resouceGroup.Do(resouceKey, func() (interface{}, error) {
		subject, err := s.cache.GetSubject(context.Background(), uint64(in.ID))
		if err == nil {
			return subject, nil
		}
		if err != ErrCacheMiss {
			return nil, err
		}

		// load From DB
		esubject, err := s.repo.GetSubject(context.Background(), uint64(in.ID))
		if err != nil {
			return nil, err
		}
		subject = newSubjectFromEntity(&esubject)
		if err := s.cache.StoreSubject(context.Background(), subject); err != nil {
			flog.Error(err, "s.cache.StoreSubject()", "subject", subject)
		}

		return subject, nil
	})

	if err != nil {
		if err == domain.DBRecordNotFound {
			return &app.GetSubjectResp{Status: apicode.ErrReplySvcNotFound.ToStatus()}, nil
		}
		return &app.GetSubjectResp{Status: apicode.ErrReplySvcGetSubject.ToStatus()}, nil
	}

	subject := res.(*dto.ReplyCommentSubject)
	if shared {
		subject = proto.Clone(subject).(*dto.ReplyCommentSubject)
	}

	return &app.GetSubjectResp{Status: apicode.StatusOK.ToStatus(), Subject: subject}, nil
}

func (s *ReplyCommentService) sortCommentToSlice(m map[uint64]*dto.ReplyComment) ([]*dto.ReplyComment, error) {
	comments := make([]*dto.ReplyComment, 0, len(m))

	for _, v := range m {
		if v.Item == nil {
			continue
		}
		comments = append(comments, v)
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Item.Floor < comments[j].Item.Floor
	})

	return comments, nil
}

func (s *ReplyCommentService) loadCommentContentToComments(ctx context.Context, ids []uint64, comments map[uint64]*dto.ReplyComment) error {
	missIds := []uint64{}
	for _, idx := range ids {
		if _, exist := comments[idx]; !exist {
			comments[idx] = &dto.ReplyComment{}
		}

		content, err := s.cache.GetCommentContent(ctx, int64(idx))
		if err == ErrCacheMiss {
			missIds = append(missIds, idx)
			continue
		} else if err != nil {
			flog.Error(err, "cache.GetCommentContent() err:")
			continue
		}

		comments[idx].Content = content
	}

	if 0 < len(missIds) {
		idx, err := s.repo.ListCommetContents(ctx, missIds)
		if err != nil {
			flog.Error(err, "loadCommentContents() err")
		}

		for _, v := range idx {
			d := &dto.ReplyCommentContent{
				ID:       int64(v.Id),
				Content:  []byte(v.Content),
				IP:       v.Ip,
				Platform: int64(v.Platform),
				Device:   v.Device,
				State:    int64(v.State),
			}
			comments[v.Id].Content = d
			if err := s.cache.StoreCommentContent(ctx, d); err != nil {
				flog.Error(err, "StoreCommentIndex()", "index", d)
			}
		}
	}

	return nil
}

func (s *ReplyCommentService) loadCommentItemToComments(ctx context.Context, ids []uint64, comments map[uint64]*dto.ReplyComment) error {
	missIds := []uint64{}
	for _, idx := range ids {
		if _, exist := comments[idx]; !exist {
			comments[idx] = &dto.ReplyComment{}
		}

		index, err := s.cache.GetCommentItem(ctx, idx)
		if err == ErrCacheMiss {
			missIds = append(missIds, idx)
			continue
		} else if err != nil {
			flog.Error(err, "cache.GetCommentItem() err:")
			continue
		}

		comments[idx].Item = index
	}

	// load miss's comment
	if 0 < len(missIds) {
		idx, err := s.repo.ListCommetItem(ctx, missIds)
		if err != nil {
			flog.Error(err, "loadCommentContents() err")
		}

		for _, v := range idx {
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
				Seq:      int64(v.Seq),
			}
			comments[v.Id].Item = d
			if err := s.cache.StoreCommentItem(ctx, d); err != nil {
				flog.Error(err, "StoreCommentItem()", "index", d)
			}
		}
	}

	return nil
}

const numPerPage = 10

func pageFromFloor(floor int64) int64 {
	page := floor / int64(numPerPage)
	return page
}

func (s *ReplyCommentService) getCommentPageIds(ctx context.Context, subject, parent, page, floor uint64) ([]uint64, error) {
	// get page from cache
	ids, err := s.cache.GetCommentPageIds(ctx, subject, parent, page)
	if err == ErrCacheMiss {
		// load index data from db to cache
		// if page hasn't data(cause by delete), skip to next page.
		entitys, err := s.repo.ListCommentsPageIds(ctx, subject, parent, int(floor), numPerPage)
		if err != nil {
			return nil, err
		}
		if len(entitys) == 0 {
			return nil, nil
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

func (s *ReplyCommentService) getSubjectLastFloor(ctx context.Context, subject, parent uint64) (int64, error) {
	last, err := s.cache.GetLastFloor(ctx, subject, parent)
	if err != nil && err != ErrCacheMiss {
		return -1, err
	}

	// key not exist
	if err == ErrCacheMiss {
		// rebuild last floor from db
		key := fmt.Sprintf(lastFloorGroupFromat, subject, parent)

		_, err, _ = s.resouceGroup.Do(key, func() (interface{}, error) {
			last, err := s.repo.GetCommentLastFloor(context.Background(), subject, parent)
			if err != nil {
				return 0, err
			}
			// store last floor to cache
			if err := s.cache.StoretLastFloorIfNotExist(context.Background(), subject, parent, last); err != nil {
				return 0, err
			}

			return last, err
		})

		if err != nil {
			return 0, err
		}

		// last floor number would be changed in the situation of concurrence
		// retry after rebuild last floor
		last, err = s.cache.GetLastFloor(ctx, subject, parent)
		if err != nil {
			return 0, err
		}
	}

	return last, err
}

// func (s *ReplyCommentService) incrCommentReplyCnt(ctx context.Context, id, subject, parent, floor uint64) error {
// 	err := s.cache.IncrCommentReplyCnt(ctx, id)
// 	if err != nil && err == ErrCacheMiss {
// 		// load page
// 		s.ListCommentPage(ctx, &app.ListCommentPageReq{Subject: int64(subject), Parent: int64(parent), Floor: int64(floor)})
// 		// try again
// 		err = s.cache.IncrCommentReplyCnt(ctx, id)
// 		if err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *ReplyCommentService) incrSubjectReplyCnt(ctx context.Context, id uint64) error {
// 	err := s.cache.IncrSubjectReplyCnt(ctx, id)
// 	if err != nil && err == ErrCacheMiss {
// 		// load cache
// 		s.GetSubject(ctx, &app.GetSubjectReq{ID: int64(id)})
// 		// try again
// 		err = s.cache.IncrSubjectReplyCnt(ctx, id)
// 		if err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}
// 	return nil

// }
