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
	CacheKeyCommentItem     = "micro:reply-svc:comment:item:%d"

	CKCommentPage            = "micro:reply-svc:comment:page:%d:%d:%d" // subject/parent/page
	CKCommentItemPage        = "micro:reply-svc:comment:item:page:%d:%d:%d"
	CKCommentContentPage     = "micro:reply-svc:comment:content:page:%d:%d:%d"
	CKLastFloor              = "micro:reply-svc:last:floor:%d:%d" // subject/parent
	CKCommentItemAttr        = "micro:reply-svc:comment:item:attr:%d"
	CKCommentItemSubjectAttr = "micro:reply-svc:subject:attr:%d"
	CKCommentsubject         = "micro:reply-svc:subject:%d"

	IndexAttrSub      = "sub"
	IndexAttrParent   = "parent"
	IndexAttrFloor    = "floor"
	IndexAttrUid      = "uid"
	IndexAttrReplyto  = "to"
	IndexAttrLike     = "like"
	IndexAttrDislike  = "dislike"
	IndexAttrReplyCnt = "reply-cnt"
	IndexAttrState    = "state"
	IndexAttrSeq      = "seq"
	AttrObjID         = "obj-id"
	AttrObjType       = "obj-type"

// IndexAttrContentID = "cid"
)

var (
	ErrCacheMiss = errors.New("Miss Cache")
)

type bytes []byte

// TODO improve by local cache
type CommentCache struct {
	rdb            *redis.Client
	cacheDura      time.Duration
	indexcacheDura time.Duration
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

	return nil
}

func (c *CommentCache) GetSubject(ctx context.Context, id uint64) (*dto.ReplyCommentSubject, error) {
	key := fmt.Sprintf(CKCommentsubject, id)
	d := &dto.ReplyCommentSubject{}
	res := c.rdb.HGetAll(ctx, key)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	attrs, err := res.Result()
	if err != nil {
		return nil, err
	}

	var val int
	d.ID = int64(id)

	val, _ = strconv.Atoi(attrs[AttrObjID])
	d.ObjID = int64(val)

	val, _ = strconv.Atoi(attrs[AttrObjType])
	d.ObjType = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrLike])
	d.Like = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrDislike])
	d.Dislike = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrReplyCnt])
	d.ReplyCnt = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrSeq])
	d.Seq = int64(val)

	return d, nil
}

func (c *CommentCache) StoreSubject(ctx context.Context, subject *dto.ReplyCommentSubject) error {
	key := fmt.Sprintf(CKCommentsubject, subject.ID)
	return c.StoreAttr(ctx, key, 6,
		[]string{AttrObjID, AttrObjType, IndexAttrLike, IndexAttrDislike, IndexAttrReplyCnt, IndexAttrSeq},
		[]interface{}{subject.ObjID, subject.ObjType, subject.Like, subject.Dislike, subject.ReplyCnt, subject.Seq})
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

	return nil
}

func (c *CommentCache) GetCommentItem(ctx context.Context, id uint64) (*dto.ReplyCommentItem, error) {
	key := fmt.Sprintf(CacheKeyCommentItem, id)
	d := &dto.ReplyCommentItem{}
	res := c.rdb.HGetAll(ctx, key)
	if res.Err() == redis.Nil {
		return nil, ErrCacheMiss
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	attrs, err := res.Result()
	if err != nil {
		return nil, err
	}

	var val int
	d.ID = int64(id)

	val, _ = strconv.Atoi(attrs[IndexAttrSub])
	d.Subject = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrParent])
	d.Parent = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrFloor])
	d.Floor = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrUid])
	d.UserID = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrReplyto])
	d.Replyto = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrLike])
	d.Like = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrDislike])
	d.Dislike = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrReplyCnt])
	d.ReplyCnt = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrState])
	d.State = int64(val)

	val, _ = strconv.Atoi(attrs[IndexAttrSeq])
	d.Seq = int64(val)

	return d, nil
}

func (c *CommentCache) StoreCommentIndex(ctx context.Context, idx *dto.ReplyCommentItem) error {
	key := fmt.Sprintf(CacheKeyCommentItem, idx.ID)
	var attrKV []interface{}
	attrKV = append(attrKV,
		IndexAttrSub, idx.Subject,
		IndexAttrParent, idx.Parent,
		IndexAttrFloor, idx.Floor,
		IndexAttrUid, idx.UserID,
		IndexAttrReplyto, idx.Replyto,
		IndexAttrLike, idx.Like,
		IndexAttrDislike, idx.Dislike,
		IndexAttrReplyCnt, idx.ReplyCnt,
		IndexAttrState, idx.State,
		IndexAttrSeq, idx.Seq,
		//IndexAttrContentID, idx.ContentID
	)

	res := c.rdb.HSet(ctx, key, attrKV...)
	if res.Err() != nil {
		return res.Err()
	}

	res2 := c.rdb.Expire(ctx, key, c.cacheDura)
	if res2.Err() != nil {
		flog.Warn("rdb.Expire() err", "err", res2.Err())
	}

	return nil
}

func (c *CommentCache) GetCommentItemPage(ctx context.Context, subject, parent, page uint64) ([]*dto.ReplyCommentItem, error) {
	key := fmt.Sprintf(CKCommentItemPage, subject, parent, page)

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
	key := fmt.Sprintf(CKCommentContentPage, subject, parent, page)
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

// const scriptStoreItemsToHashIfNotExist = `
//   local itemsKey = KEYS[1]

//   local key_exists = redis.call('EXISTS', itemsKey)
//   if key_exists == 1 then
// 	return 0
//   end

//   local res = redis.call('HSET', itemsKey, unpack(ARGV))

//   return 1
// `

// func (c *CommentCache) SoreCommentItemPage(ctx context.Context, subject, parent, page uint64, items []*dto.ReplyCommentItem) error {
// 	key := fmt.Sprintf(CKCommentItemPage, subject, parent, page)

// 	itemLen := len(items)
// 	scriptArgs := make([]string, 0, itemLen*2)

// 	for _, v := range items {
// 		i, err := proto.Marshal(v)
// 		if err != nil {
// 			flog.Error(err, "proto.Marshal(v)", "v", v)
// 			continue
// 		}

// 		scriptArgs = append(scriptArgs, strconv.Itoa(int(v.ID)))
// 		scriptArgs = append(scriptArgs, string(i))
// 	}

// 	res := c.rdb.Eval(ctx, scriptStoreItemsToHashIfNotExist, []string{key}, scriptArgs)
// 	if res.Err() != nil {
// 		return res.Err()
// 	}

// 	for _, v := range items {
// 		c.StoreCommentItemAttr(ctx, v)
// 	}

// 	state, err := res.Int()
// 	if err != nil {
// 		return err
// 	}

// 	if state == 0 {
// 		// key exists, not thing to do
// 		return nil
// 	}

// 	c.rdb.Expire(ctx, key, c.cacheDura)

// 	return nil
// }

// func (c *CommentCache) SoreCommentContentPage(ctx context.Context, subject, parent, page uint64, items []*dto.ReplyCommentContent) error {
// 	key := fmt.Sprintf(CKCommentItemPage, subject, parent, page)

// 	itemLen := len(items)
// 	scriptArgs := make([]string, 0, itemLen*2)

// 	for _, v := range items {
// 		i, err := proto.Marshal(v)
// 		if err != nil {
// 			flog.Error(err, "proto.Marshal(v)", "v", v)
// 			continue
// 		}

// 		scriptArgs = append(scriptArgs, strconv.Itoa(int(v.ID)))
// 		scriptArgs = append(scriptArgs, string(i))
// 	}

// 	res := c.rdb.Eval(ctx, scriptStoreItemsToHashIfNotExist, []string{key}, scriptArgs)
// 	if res.Err() != nil {
// 		return res.Err()
// 	}

// 	state, err := res.Int()
// 	if err != nil {
// 		return err
// 	}

// 	if state == 0 {
// 		// key exists, not thing to do
// 		return nil
// 	}

// 	c.rdb.Expire(ctx, key, c.cacheDura)

// 	return nil
// }

const scriptIncrIfExist = `
local key = KEYS[1]

local key_exists = redis.call('EXISTS', itemsKey)
if key_exists == 0 then
  return -1
end

local res = redis.call('incr', key)

return res
`

func (c *CommentCache) GetLastFloor(ctx context.Context, subject, parent uint64) (int64, error) {
	key := fmt.Sprintf(CKLastFloor, subject, parent)
	res := c.rdb.Eval(ctx, scriptIncrIfExist, []string{key}, nil)
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

local key_exists = redis.call('EXISTS', itemsKey)
if key_exists == 1 then
  return 0
end

local res = redis.call('set', key, ARGV[1] EX ARGV[2])

return res
`

func (c *CommentCache) StoretLastFloorIfNotExist(ctx context.Context, subject, parent uint64, floor uint64) error {
	key := fmt.Sprintf(CKLastFloor, subject, parent)
	res := c.rdb.Eval(ctx, scriptSetIfNotExist, []string{key}, []interface{}{floor, c.cacheDura})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// const scriptPutCommentItem = `
//   local page_key = KEYS[1]
//   local attr_key = KEYS[2]
//   local item_key = KEYS[3]

//   local seq = ARGV[1]
//   local like = ARGV[2]
//   local dislike = ARGV[3]
//   local reply_cnt = ARGV[4]
//   local item = ARGV[5]

//   local res = redis.call('hget', attr_key, 'seq')
//   if res ~= nil and seq < res then
//     return 1
//   end

//   res = redis.call('hget', attr_key, 'seq', seq, 'like', like, 'dislike',dislike, 'reply_cnt',reply_cnt)
//   if res ~= 1 then
//     return res
//   end

//   res = redis.call('EXISTS', page_key)
//   if res == 0 then
//     return 0
//   end

//   res = redis.call('hget', page_key, item_key, item)
//   if res ~= 1 then
//     return res
//   end

//   return 1
// `

// func (c *CommentCache) PutCommentItem(ctx context.Context, item *dto.ReplyCommentItem) error {
// 	like := item.Like
// 	dislike := item.Dislike
// 	replyCnt := item.ReplyCnt
// 	seq := item.Seq

// 	page := pageFromFloor(item.Floor)
// 	pageKey := fmt.Sprintf(CKCommentItemPage, item.Subject, item.Subject, page)
// 	attrKey := fmt.Sprintf(CKCommentItemAttr, item.ID)
// 	itemkey := strconv.Itoa(int(item.ID))

// 	itembytes, err := proto.Marshal(item)
// 	if err != nil {
// 		return err
// 	}

// 	res := c.rdb.Eval(ctx, scriptSetIfNotExist, []string{pageKey, attrKey, itemkey}, []interface{}{seq, like, dislike, replyCnt, itembytes})

// 	if res.Err() != nil {
// 		return res.Err()
// 	}
// 	state, err := res.Int()
// 	if err != nil {
// 		return err
// 	}

// 	if state == 0 {
// 		return ErrCacheMiss
// 	}

// 	return nil

// }

const scriptStoreCommentItemAttr = `
local attr_key = KEYS[1]

local kvlen = ARGV[1]

local seq = ARGV[2]

local res = redis.call('hget', attr_key, 'seq')

if res ~= nil and res ~= false and seq <= res then 
  return 1
end

local attr_kv = {}

local j = 2
for i = 1, kvlen*2, 2
do
  attr_kv[i] = KEYS[j]
  attr_kv[i+1] = ARGV[j]
  j = j+1
end

res = redis.call('hset', attr_key, unpack(attr_kv))
if res ~= 1 then 
  return res
end

return 1
`

// require seq filed
func (c *CommentCache) StoreAttr(ctx context.Context, attrKey string, len int, attrs []string, vals []interface{}) error {
	res := c.rdb.Eval(ctx, scriptStoreCommentItemAttr, append([]string{attrKey}, attrs...), append([]interface{}{len}, vals...))
	return res.Err()
}

func (c *CommentCache) StoreCommentItemAttr(ctx context.Context, item *dto.ReplyCommentItem) error {
	like := item.Like
	dislike := item.Dislike
	replyCnt := item.ReplyCnt
	seq := item.Seq

	attrKey := fmt.Sprintf(CKCommentItemAttr, item.ID)

	// res := c.rdb.Eval(ctx, scriptStoreCommentItemAttr,
	// 	[]string{attrKey, IndexAttrSeq, IndexAttrLike, IndexAttrDislike, IndexAttrReplyCnt, IndexAttrSub, IndexAttrParent, IndexAttrFloor, IndexAttrReplyto, IndexAttrState, IndexAttrUid},
	// 	[]interface{}{10, seq, like, dislike, replyCnt, item.Subject, item.Parent, item.Floor, item.Replyto, item.State, item.ID})

	return c.StoreAttr(ctx, attrKey, 10,
		[]string{IndexAttrSeq, IndexAttrLike, IndexAttrDislike, IndexAttrReplyCnt, IndexAttrSub, IndexAttrParent, IndexAttrFloor, IndexAttrReplyto, IndexAttrState, IndexAttrUid},
		[]interface{}{seq, like, dislike, replyCnt, item.Subject, item.Parent, item.Floor, item.Replyto, item.State, item.ID})

}

// func (c *CommentCache) StoreCommentContentIfExist(ctx context.Context, subject, parent, page uint64, content *dto.ReplyCommentContent) error {
// 	pageKey := fmt.Sprintf(CKCommentContentPage, subject, parent, page)

// 	contenBytes, err := proto.Marshal(content)
// 	if err != nil {
// 		return err
// 	}
// 	res := c.rdb.Exists(ctx, pageKey)
// 	if state, err := res.Result(); err != nil {
// 		return err
// 	} else if state == 0 {
// 		// page key not exist. it is not necessary to cache this content.
// 		return nil
// 	}

// 	res = c.rdb.HSet(ctx, pageKey, content.ID, contenBytes)
// 	if res.Err() != nil {
// 		return err
// 	}

// 	return nil
// }

const scriptHIncrWithSeq = `
local key = KEYS[1]
local field = KEYS[2]

local key_exists = redis.call('EXISTS', attr_key)
if key_exists == 1 then
  return 0
end


redis.call('HINCRBY', key, field, 1)

redis.call('HINCRBY', key, "seq", 1)

return 1
`

func (c *CommentCache) IncrCommentReplyCnt(ctx context.Context, id uint64) error {
	key := fmt.Sprintf(CKCommentItemAttr, id)

	res := c.rdb.Eval(ctx, scriptHIncrWithSeq, []string{key, IndexAttrReplyCnt})

	if state, err := res.Int(); err != nil {
		return err
	} else if state == 0 {
		return ErrCacheMiss
	}

	return res.Err()
}

func (c *CommentCache) IncrSubjectReplyCnt(ctx context.Context, id uint64) error {
	key := fmt.Sprintf(CKCommentItemSubjectAttr, id)

	res := c.rdb.Eval(ctx, scriptHIncrWithSeq, []string{key, IndexAttrReplyCnt})

	if state, err := res.Int(); err != nil {
		return err
	} else if state == 0 {
		return ErrCacheMiss
	}

	return res.Err()
}
