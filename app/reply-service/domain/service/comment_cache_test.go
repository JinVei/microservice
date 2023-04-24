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
		rdb:           rediscli,
		cacheDra:      time.Hour,
		indexcacheDra: time.Hour,
	}
	testdata := []*dto.ReplyCommentItem{
		{
			ID:    1,
			Count: 11,
		},
		{
			ID:    2,
			Count: 12,
		},
		{
			ID:    3,
			Count: 13,
		},
		{
			ID:    4,
			Count: 14,
		},
		{
			ID:    5,
			Count: 15,
		},
	}

	err := cm.SoreCommentIndexPage(context.Background(), 0, 1, 0, testdata)
	if err != nil {
		t.Fatal(err)
	}

}
