package service

import (
	"context"
	"time"

	"github.com/jinvei/microservice/app/reply-service/domain"
	"github.com/jinvei/microservice/app/reply-service/domain/entity"
	"github.com/jinvei/microservice/base/api/proto/v1/dto"
)

const (
	CountableBatchItemSubject = 0
	CountableBatchItemComment = 0
)

type CommentBatchWriter struct {
	repo  domain.IReplyCommentRepository
	cache *CommentCache
	//	ctx  context.Context

	PutCommentC chan *dto.ReplyComment
	SubjectCntC chan entity.CountableItem
	CommentCntC chan entity.CountableItem

	BatchNum int
	Interval time.Duration
}

// TODO: how to deal with like a comment that not summit.
func NewCommentBatchWriter(repo domain.IReplyCommentRepository, c *CommentCache, BatchNum int, Interval time.Duration) *CommentBatchWriter {
	return &CommentBatchWriter{
		repo:        repo,
		BatchNum:    BatchNum,
		Interval:    Interval,
		PutCommentC: make(chan *dto.ReplyComment, BatchNum/2),
		SubjectCntC: make(chan entity.CountableItem, BatchNum/2),
		CommentCntC: make(chan entity.CountableItem, BatchNum/2),
		cache:       c,
	}
}

func (w *CommentBatchWriter) Start() {
	go w._batchPutCommentCoroutine()
	go w._batchIncrCommentCoroutine()
	go w._batchIncrSubjectCoroutine()
}

func (w *CommentBatchWriter) Close() {
	close(w.PutCommentC)
	close(w.SubjectCntC)
	close(w.CommentCntC)
}

func (w *CommentBatchWriter) PutComment(comment *dto.ReplyComment) error {
	w.PutCommentC <- comment
	return nil
}

func (w *CommentBatchWriter) InrcCommentLikeCnt(id uint64) error {
	w.CommentCntC <- entity.CountableItem{Id: id, Like: 1}
	return nil
}

func (w *CommentBatchWriter) InrcCommentReplyCnt(id uint64) error {
	w.CommentCntC <- entity.CountableItem{Id: id, Reply: 1}
	return nil
}

func (w *CommentBatchWriter) InrcSubjectReplyCnt(id uint64) error {
	w.SubjectCntC <- entity.CountableItem{Id: id, Reply: 1}
	return nil
}

func (w *CommentBatchWriter) InrcSubjectLikeCnt(id uint64) error {
	w.SubjectCntC <- entity.CountableItem{Id: id, Like: 1}
	return nil
}

func (w *CommentBatchWriter) _batchPutCommentCoroutine() {
	flog.Info("Start batch Put Comment Coroutine", "BatchNum", w.BatchNum, "Ticker", w.Interval.Seconds())
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	batch := make([]*dto.ReplyComment, 0, w.BatchNum)
	exit := false
	for {
		if exit {
			break
		}

		select {
		case c, ok := <-w.PutCommentC:
			if !ok {
				exit = true
				break
			}

			batch = append(batch, c)
			if len(batch) < w.BatchNum {
				continue
			}

			batch = w.batchSubmitComment(batch)
			ticker.Reset(w.Interval)

		case <-ticker.C:
			batch = w.batchSubmitComment(batch)
			ticker.Reset(w.Interval)
		}
	}
	if 0 < len(batch) {
		w.batchSubmitComment(batch)
	}
}

func (w *CommentBatchWriter) _batchIncrSubjectCoroutine() {
	flog.Info("Start batch Incr Subject Count Coroutine", "BatchNum", w.BatchNum, "Ticker", w.Interval.Seconds())
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	batchs := make(map[uint64]entity.CountableItem, w.BatchNum)
	exit := false

	for {
		if exit {
			break
		}

		select {
		case c, ok := <-w.SubjectCntC:
			if !ok {
				exit = true
				break
			}

			if it, exist := batchs[c.Id]; !exist {
				batchs[c.Id] = c
			} else {
				it.Like += c.Like
				it.Reply += c.Reply
				batchs[c.Id] = it
			}

			if len(batchs) < w.BatchNum {
				continue
			}

			s := cntMoveToSlice(batchs)
			w.batchIncrSubjectCnt(s)

			ticker.Reset(w.Interval)

		case <-ticker.C:
			if len(batchs) <= 0 {
				continue
			}
			s := cntMoveToSlice(batchs)
			w.batchIncrSubjectCnt(s)

			ticker.Reset(w.Interval)
		}
	}
	if 0 < len(batchs) {
		s := cntMoveToSlice(batchs)
		w.batchIncrSubjectCnt(s)
	}
}

func (w *CommentBatchWriter) _batchIncrCommentCoroutine() {
	flog.Info("Start batch Incr Comment Count Coroutine", "BatchNum", w.BatchNum, "Ticker", w.Interval.Seconds())
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	batchs := make(map[uint64]entity.CountableItem, w.BatchNum)
	exit := false

	for {
		if exit {
			break
		}

		select {
		case c, ok := <-w.CommentCntC:
			if !ok {
				exit = true
				break
			}

			if it, exist := batchs[c.Id]; !exist {
				batchs[c.Id] = c
			} else {
				it.Like += c.Like
				it.Reply += c.Reply
				batchs[c.Id] = it
			}

			if len(batchs) < w.BatchNum {
				continue
			}

			s := cntMoveToSlice(batchs)
			w.batchIncrCommentCnt(s)

			ticker.Reset(w.Interval)

		case <-ticker.C:
			if len(batchs) <= 0 {
				continue
			}
			s := cntMoveToSlice(batchs)
			w.batchIncrCommentCnt(s)

			ticker.Reset(w.Interval)
		}
	}
	if 0 < len(batchs) {
		s := cntMoveToSlice(batchs)
		w.batchIncrCommentCnt(s)
	}
}

func (w *CommentBatchWriter) batchSubmitComment(batchs []*dto.ReplyComment) []*dto.ReplyComment {
	if len(batchs) <= 0 {
		return batchs
	}

	flog.Info("Batch Submit Comment", "batchs", batchs)

	items, contents, err := w.repo.BatchSubmitComments(context.Background(), batchs)
	if err != nil {
		flog.Error(err, "w.repo.BatchSubmitComments() error", "batch", batchs)
		// TODO: send err to event bus
	}
	batchs = batchs[0:0]

	for _, item := range items {
		page := pageFromFloor(int64(item.Floor))
		if err := w.cache.StoreCommentPageIds(context.Background(), uint64(item.Subject), uint64(item.Parent), uint64(page), []uint64{uint64(item.Id)}); err != nil {
			flog.Error(err, "w.repo.StoreCommentPageIds() error", "item", item)
		}

		if err := w.cache.StoreCommentItem(context.Background(), newCommentItemFromEntity(&item)); err != nil {
			flog.Error(err, "cache.StoreCommentItem() err", "item", item)
		}

		if err := w.InrcSubjectReplyCnt(item.Subject); err != nil {
			flog.Error(err, "w.InrcSubjectReplyCnt err", "id", item.Id)
		}

		if item.Parent != 0 {
			if err := w.InrcCommentReplyCnt(item.Parent); err != nil {
				flog.Error(err, "w.InrcCommentReplyCnt err", "id", item.Id)
			}
		}
	}

	for _, content := range contents {
		if err := w.cache.StoreCommentContent(context.Background(), newCommentContentFromEntity(&content)); err != nil {
			flog.Error(err, "cache.StoreCommentContent() err", "content", content)
		}
	}

	return batchs
}

func (w *CommentBatchWriter) batchIncrCommentCnt(batchs []entity.CountableItem) {
	if len(batchs) <= 0 {
		return
	}

	flog.Info("Batch Incr Comment Cnt", "batchs", batchs)
	items, contents, err := w.repo.BatchIncrCommentCount(context.Background(), batchs)
	if err != nil {
		flog.Error(err, "w.repo.batchIncrCommentCnt() error", "batch", batchs)
		// TODO: send err to event bus
	}

	for _, item := range items {
		if err := w.cache.StoreCommentItem(context.Background(), newCommentItemFromEntity(item)); err != nil {
			flog.Error(err, "w.cache.StoreCommentItem(context.Background(), newCommentItemFromEntity(&i)) error", "item", item)
		}
	}

	for _, c := range contents {
		if err := w.cache.StoreCommentContent(context.Background(), newCommentContentFromEntity(c)); err != nil {
			flog.Error(err, "w.cache.StoreCommentContent(context.Background(), newCommentContentFromEntity(&c)) error", "content", c)
		}
	}
}

func (w *CommentBatchWriter) batchIncrSubjectCnt(batchs []entity.CountableItem) {
	if len(batchs) <= 0 {
		return
	}

	flog.Info("Batch Incr subject Cnt", "batchs", batchs)
	subjects, err := w.repo.BatchIncrSubjectCount(context.Background(), batchs)
	if err != nil {
		flog.Error(err, "w.repo.batchIncrSubjectCnt() error", "batch", batchs)
		// TODO: send err to event bus
	}

	for _, subj := range subjects {
		if err := w.cache.StoreSubject(context.Background(), newSubjectFromEntity(&subj)); err != nil {
			flog.Error(err, "w.cache.StoreSubject(context.Background(), newSubjectFromEntity(&subj)) error", "item", subj)
		}
	}
}

func cntMoveToSlice(m map[uint64]entity.CountableItem) (s []entity.CountableItem) {
	for k, v := range m {
		s = append(s, v)
		delete(m, k)
	}
	return s
}
