package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/jinvei/microservice/base/framework/cache"
	"github.com/jinvei/microservice/base/framework/configuration"
)

func TestCacheStoreCommentItemsPage(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	conf := configuration.DefaultOrDie()
	conf.SetSystemID("10001")

	rediscli := cache.RedisClient(conf)

	if rediscli == nil {
		t.Fatal()
	}

	cm := CommentCache{
		rdb:       rediscli,
		cacheDura: time.Hour,
		indexDura: time.Hour,
	}
	testdata := []*dto.ReplyCommentItem{
		{
			ID:       1,
			ReplyCnt: 11,
		},
		{
			ID:       2,
			ReplyCnt: 12,
		},
		{
			ID:       3,
			ReplyCnt: 13,
		},
		{
			ID:       4,
			ReplyCnt: 14,
		},
		{
			ID:       5,
			ReplyCnt: 15,
		},
	}

	// err := cm.SoreCommentItemPage(context.Background(), 0, 1, 0, testdata)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	err := cm.StoreCommentItemAttr(context.Background(), testdata[0])
	if err != nil {
		t.Fatal(err)
	}

}

func TestCacheKeyNotExist(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	conf := configuration.DefaultOrDie()
	conf.SetSystemID("10001")

	rediscli := cache.RedisClient(conf)
	res := rediscli.Get(context.Background(), "111")
	if res.Err() != nil {
		t.Fatal(res.Err())
	}
}
