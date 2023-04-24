package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const (
	CacheKeyCommentContennt = "micro:reply-svc:comment:content:%d"
	CacheKeyCommentIndex    = "micro:reply-svc:comment:index:%d"

	CacheKeyCommentPage        = "micro:reply-svc:comment:page:%d:%d:%d" // subject/parent/page
	CacheKeyCommentIndexPage   = "micro:reply-svc:comment:index:page:%d:%d:%d"
	CacheKeyCommentContentPage = "micro:reply-svc:comment:content:page:%d:%d:%d"

	IndexAttrSub       = "sub"
	IndexAttrParent    = "parent"
	IndexAttrFloor     = "floor"
	IndexAttrUid       = "uid"
	IndexAttrReplyto   = "to"
	IndexAttrLike      = "like"
	IndexAttrHate      = "hate"
	IndexAttrCount     = "cnt"
	IndexAttrState     = "state"
	IndexAttrContentID = "cid"
)

var (
	ErrCacheMiss = errors.New("Miss Cache")
)

type bytes []byte

// TODO improve by local cache
type CommentCache struct {
	rdb           *redis.Client
	cacheDra      time.Duration
	indexcacheDra time.Duration
}

// GetCommentIndex get Comment Index from cache
func (c *CommentCache) GetCommentPageIds(ctx context.Context, subject, parent, page int64) ([]uint64, error) {
	key := fmt.Sprintf(CacheKeyCommentPage, subject, parent, page)

	res := c.rdb.ZRevRange(ctx, key, 0, -1)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	d, err := res.Result()
	if err != nil {
		return nil, err
	}

	indexes := make([]uint64, 0, len(d))
	for _, v := range d {
		i, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		indexes = append(indexes, uint64(i))
	}

	return indexes, nil
}

// func (c *CommentCache) StoreCommentPageIds(ctx context.Context, subject, parent, page int64, ids []uint64) error {
// 	key := fmt.Sprintf(CacheKeyCommentPage, subject, parent, page)
// 	c.rdb.ZAdd(key, ids...)
// }

// func (c *CommentCache) GetCommentContent(ctx context.Context, id int64) (*dto.ReplyCommentContent, error) {
// 	key := fmt.Sprintf(CacheKeyCommentContennt, id)
// 	d := &dto.ReplyCommentContent{}
// 	res := c.rdb.Get(ctx, key)

// 	if res.Err() == redis.Nil {
// 		return nil, ErrCacheMiss
// 	}
// 	if res.Err() != nil {
// 		return nil, res.Err()
// 	}
// 	b, err := res.Bytes()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := proto.Unmarshal(b, d); err != nil {
// 		return nil, err
// 	}
// 	return d, nil
// }

// func (c *CommentCache) StoreCommentContent(ctx context.Context, content *dto.ReplyCommentContent) error {
// 	key := fmt.Sprintf(CacheKeyCommentContennt, content.ID)
// 	b, err := proto.Marshal(content)
// 	if err != nil {
// 		return err
// 	}
// 	status := c.rdb.Set(ctx, key, b, c.cacheDra)
// 	if status.Err() != nil {
// 		return status.Err()
// 	}

// 	return nil
// }

func (c *CommentCache) StoreCommentPageIds(ctx context.Context, subject, parent, page int64, ids []uint64) error {
	key := fmt.Sprintf(CacheKeyCommentPage, subject, parent, page)
	members := make([]redis.Z, 0, len(ids))

	for i, v := range ids {
		members = append(members, redis.Z{Member: v, Score: float64(i)})
	}
	res := c.rdb.ZAdd(ctx, key, members...)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// func (c *CommentCache) GetCommentIndex(ctx context.Context, id uint64) (*dto.ReplyCommentItem, error) {
// 	key := fmt.Sprintf(CacheKeyCommentIndex, id)
// 	d := &dto.ReplyCommentItem{}
// 	res := c.rdb.HGetAll(ctx, key)
// 	if res.Err() == redis.Nil {
// 		return nil, ErrCacheMiss
// 	}
// 	if res.Err() != nil {
// 		return nil, res.Err()
// 	}

// 	attrs, err := res.Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var val int
// 	d.ID = int64(id)
// 	val, _ = strconv.Atoi(attrs[IndexAttrSub])
// 	d.Subject = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrParent])
// 	d.Parent = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrFloor])
// 	d.Floor = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrUid])
// 	d.UserID = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrReplyto])
// 	d.Replyto = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrLike])
// 	d.Like = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrHate])
// 	d.Hate = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrCount])
// 	d.Count = int64(val)

// 	val, _ = strconv.Atoi(attrs[IndexAttrState])
// 	d.State = int64(val)

// 	// val, _ = strconv.Atoi(attrs[IndexAttrContentID])
// 	// d.ContentID = int64(val)

// 	//attrs[]
// 	return d, nil
// }

// func (c *CommentCache) StoreCommentIndex(ctx context.Context, idx *dto.ReplyCommentItem) error {
// 	key := fmt.Sprintf(CacheKeyCommentIndex, idx.ID)
// 	var attrKV []interface{}
// 	attrKV = append(attrKV,
// 		IndexAttrSub, idx.Subject,
// 		IndexAttrParent, idx.Parent,
// 		IndexAttrFloor, idx.Floor,
// 		IndexAttrUid, idx.UserID,
// 		IndexAttrReplyto, idx.Replyto,
// 		IndexAttrLike, idx.Like,
// 		IndexAttrHate, idx.Hate,
// 		IndexAttrCount, idx.Count,
// 		IndexAttrState, idx.State,
// 		//IndexAttrContentID, idx.ContentID
// 	)

// 	res := c.rdb.HSet(ctx, key, attrKV...)
// 	if res.Err() != nil {
// 		return res.Err()
// 	}

// 	res2 := c.rdb.Expire(ctx, key, c.cacheDra)
// 	if res2.Err() != nil {
// 		flog.Warn("rdb.Expire() err", "err", res2.Err())
// 	}

// 	return nil
// }

func (c *CommentCache) GetCommentIndexPage(ctx context.Context, subject, parent, page uint64) ([]*dto.ReplyCommentItem, error) {
	key := fmt.Sprintf(CacheKeyCommentIndexPage, subject, parent, page)

	res := c.rdb.HGetAll(ctx, key)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	itemMap, err := res.Result()
	if err != nil {
		return nil, err
	}

	items := make([]*dto.ReplyCommentItem, 0, len(itemMap))
	for _, v := range itemMap {
		d := &dto.ReplyCommentItem{}
		if err := proto.Unmarshal([]byte(v), d); err != nil {
			flog.Error(err, "proto.Unmarshal([]byte(v), d)", "v", v)
			continue
		}
		items = append(items, d)
	}

	return items, nil
}

func (c *CommentCache) GetCommentContentPage(ctx context.Context, subject, parent, page uint64) ([]*dto.ReplyCommentContent, error) {
	key := fmt.Sprintf(CacheKeyCommentContentPage, subject, parent, page)
	res := c.rdb.HGetAll(ctx, key)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}

	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	itemMap, err := res.Result()
	if err != nil {
		return nil, err
	}

	items := make([]*dto.ReplyCommentContent, 0, len(itemMap))
	for _, v := range itemMap {
		d := &dto.ReplyCommentContent{}
		if err := proto.Unmarshal([]byte(v), d); err != nil {
			flog.Error(err, "proto.Unmarshal([]byte(v), d)", "v", v)
			continue
		}
		items = append(items, d)
	}
	return items, nil
}

const scriptStoreItemsToHashIfNotExist = `
  local itemsKey = KEYS[1]

  local key_exists = redis.call('EXISTS', itemsKey)
  if key_exists == 1 then
	return 0
  end

  local res = redis.call('HSET', itemsKey, unpack(ARGV))

  return 1
`

func (c *CommentCache) SoreCommentIndexPage(ctx context.Context, subject, parent, page uint64, items []*dto.ReplyCommentItem) error {
	key := fmt.Sprintf(CacheKeyCommentIndexPage, subject, parent, page)

	itemLen := len(items)
	scriptArgs := make([]string, 0, itemLen*2)

	for _, v := range items {
		i, err := proto.Marshal(v)
		if err != nil {
			flog.Error(err, "proto.Marshal(v)", "v", v)
			continue
		}

		scriptArgs = append(scriptArgs, strconv.Itoa(int(v.ID)))
		scriptArgs = append(scriptArgs, string(i))
	}

	res := c.rdb.Eval(ctx, scriptStoreItemsToHashIfNotExist, []string{key}, scriptArgs)
	if res.Err() != nil {
		return res.Err()
	}

	state, err := res.Int()
	if err != nil {
		return err
	}

	if state == 0 {
		// key exists, not thing to do
		return nil
	}

	c.rdb.Expire(ctx, key, c.cacheDra)

	return nil
}

func (c *CommentCache) SoreCommentContentPage(ctx context.Context, subject, parent, page uint64, items []*dto.ReplyCommentContent) error {
	key := fmt.Sprintf(CacheKeyCommentIndexPage, subject, parent, page)

	itemLen := len(items)
	scriptArgs := make([]string, 0, itemLen*2)

	for _, v := range items {
		i, err := proto.Marshal(v)
		if err != nil {
			flog.Error(err, "proto.Marshal(v)", "v", v)
			continue
		}

		scriptArgs = append(scriptArgs, strconv.Itoa(int(v.ID)))
		scriptArgs = append(scriptArgs, string(i))
	}

	res := c.rdb.Eval(ctx, scriptStoreItemsToHashIfNotExist, []string{key}, scriptArgs)
	if res.Err() != nil {
		return res.Err()
	}

	state, err := res.Int()
	if err != nil {
		return err
	}

	if state == 0 {
		// key exists, not thing to do
		return nil
	}

	c.rdb.Expire(ctx, key, c.cacheDra)

	return nil
}
