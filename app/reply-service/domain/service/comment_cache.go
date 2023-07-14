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
	// CacheKeyCommentItem     = "micro:reply-svc:comment:item:%d"
	CacheKeyCommentItem    = "comment:item:%d"
	CacheKeyData           = "data"
	CacheKeyCommentItemSeq = "seq"

	CKCommentPage        = "micro:reply-svc:comment:page:%d:%d:%d" // subject/parent/page
	CKCommentItemPage    = "micro:reply-svc:comment:item:page:%d:%d:%d"
	CKCommentContentPage = "micro:reply-svc:comment:content:page:%d:%d:%d"
	CKLastFloor          = "micro:reply-svc:last-floor:%d:%d" // subject/parent
	//CKCommentItemAttr        = "micro:reply-svc:comment:item:attr:%d"
	//CKCommentItemSubjectAttr = "micro:reply-svc:subject:attr:%d"
	CKCommentsubject = "micro:reply-svc:subject:%d"

	// IndexAttrSub      = "sub"
	// IndexAttrParent   = "parent"
	// IndexAttrFloor    = "floor"
	// IndexAttrUid      = "uid"
	// IndexAttrReplyto  = "to"
	// IndexAttrLike     = "like"
	// IndexAttrDislike  = "dislike"
	// IndexAttrReplyCnt = "reply-cnt"
	// IndexAttrState    = "state"
	AttrSeq = "seq"
	// AttrObjID         = "obj-id"
	// AttrObjType       = "obj-type"

// IndexAttrContentID = "cid"
)

var (
	ErrCacheMiss = errors.New("Miss Cache")
)

type bytes []byte

func NewCommentCache(rdb *redis.Client, cachedura, indexDura string) *CommentCache {
	cd, err := time.ParseDuration(cachedura)
	if err != nil {
		cd = time.Hour * 24 * 7 // default 7 dat
	}

	id, err := time.ParseDuration(indexDura)
	if err != nil {
		id = time.Hour * 24 * 14 // default 14 dat
	}

	return &CommentCache{
		rdb:       rdb,
		cacheDura: cd,
		indexDura: id,
	}
}

type CommentCache struct {
	rdb       *redis.Client
	cacheDura time.Duration
	indexDura time.Duration
}

// GetCommentPageIds get Comment item from cache
func (c *CommentCache) GetCommentPageIds(ctx context.Context, subject, parent, page uint64) ([]uint64, error) {
	key := fmt.Sprintf(CKCommentPage, subject, parent, page)

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

	if len(d) == 0 {
		return nil, ErrCacheMiss
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

func (c *CommentCache) GetCommentContent(ctx context.Context, id int64) (*dto.ReplyCommentContent, error) {
	key := fmt.Sprintf(CacheKeyCommentContennt, id)
	d := &dto.ReplyCommentContent{}
	res := c.rdb.Get(ctx, key)

	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	b, err := res.Bytes()
	if err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(b, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *CommentCache) StoreCommentContent(ctx context.Context, content *dto.ReplyCommentContent) error {
	key := fmt.Sprintf(CacheKeyCommentContennt, content.ID)
	b, err := proto.Marshal(content)
	if err != nil {
		return err
	}
	status := c.rdb.Set(ctx, key, b, c.cacheDura)
	if status.Err() != nil {
		return status.Err()
	}
	c.rdb.Expire(ctx, key, c.cacheDura)

	return nil
}

func (c *CommentCache) GetSubject(ctx context.Context, id uint64) (*dto.ReplyCommentSubject, error) {
	key := fmt.Sprintf(CKCommentsubject, id)
	d := &dto.ReplyCommentSubject{}
	res := c.rdb.HGet(ctx, key, CacheKeyData)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}

	data, err := res.Bytes()
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(data, d)
	return d, err
}

func (c *CommentCache) StoreSubject(ctx context.Context, subject *dto.ReplyCommentSubject) error {
	seq := subject.Seq
	key := fmt.Sprintf(CKCommentsubject, subject.ID)

	data, err := proto.Marshal(subject)
	if err != nil {
		return err
	}
	return c.StoreHashDataWithSeq(ctx, key, seq, data)
}

func (c *CommentCache) StoreCommentPageIds(ctx context.Context, subject, parent, page uint64, ids []uint64) error {
	key := fmt.Sprintf(CKCommentPage, subject, parent, page)
	members := make([]redis.Z, 0, len(ids))

	for i, v := range ids {
		members = append(members, redis.Z{Member: v, Score: float64(i)})
	}
	res := c.rdb.ZAdd(ctx, key, members...)
	if res.Err() != nil {
		return res.Err()
	}
	c.rdb.Expire(ctx, key, c.indexDura)

	return nil
}

func (c *CommentCache) GetCommentItem(ctx context.Context, id uint64) (*dto.ReplyCommentItem, error) {
	key := fmt.Sprintf(CacheKeyCommentItem, id)
	d := &dto.ReplyCommentItem{}
	res := c.rdb.HGet(ctx, key, CacheKeyData)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}

	data, err := res.Bytes()
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(data, d)
	return d, err
}

const scriptIncrIfExist = `
local key = KEYS[1]

local key_exists = redis.call('EXISTS', key)
if key_exists == 0 then
  return -1
end

local res = redis.call('incr', key)

return res
`

func (c *CommentCache) GetLastFloor(ctx context.Context, subject, parent uint64) (int64, error) {
	key := fmt.Sprintf(CKLastFloor, subject, parent)
	res := c.rdb.Eval(ctx, scriptIncrIfExist, []string{key})
	if res.Err() != nil {
		return -1, res.Err()
	}
	incr, err := res.Int64()
	if err != nil {
		return -1, err
	}
	if incr == -1 {
		return -1, ErrCacheMiss
	}

	return incr, nil
}

const scriptSetIfNotExist = `
local key = KEYS[1]

local key_exists = redis.call('EXISTS', key)
if key_exists == 1 then
  return 0
end

local res = redis.call('set', key, ARGV[1])

return res
`

func (c *CommentCache) StoretLastFloorIfNotExist(ctx context.Context, subject, parent uint64, floor uint64) error {
	key := fmt.Sprintf(CKLastFloor, subject, parent)
	res := c.rdb.Eval(ctx, scriptSetIfNotExist, []string{key}, []interface{}{floor, c.cacheDura.Seconds()})
	if res.Err() != nil {
		return res.Err()
	}
	c.rdb.Expire(ctx, key, c.indexDura)

	return nil
}

const scriptStoreHashDataWithSeq = `
local key = KEYS[1]
local seq = ARGV[1]
local data = ARGV[2]

local res = redis.call('hget', key, 'seq')

if res ~= nil and res ~= false and tonumber(seq) <= tonumber(res) then 
  return 1
end

res = redis.call('hset', key, "seq", seq, "data", data)
if res ~= 1 then 
  return res
end

return 1

`

func (c *CommentCache) StoreHashDataWithSeq(ctx context.Context, key string, seq int64, data []byte) error {
	res := c.rdb.Eval(ctx, scriptStoreHashDataWithSeq, []string{key}, seq, data)
	c.rdb.Expire(ctx, key, c.indexDura)
	return res.Err()
}

func (c *CommentCache) StoreCommentItem(ctx context.Context, item *dto.ReplyCommentItem) error {
	seq := item.Seq
	key := fmt.Sprintf(CacheKeyCommentItem, item.ID)

	data, err := proto.Marshal(item)
	if err != nil {
		return err
	}
	return c.StoreHashDataWithSeq(ctx, key, seq, data)
}
