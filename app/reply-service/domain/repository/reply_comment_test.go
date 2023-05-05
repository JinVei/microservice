package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jinvei/microservice/app/reply-service/domain/entity"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/datasource"
)

func TestCreateComment(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	conf := configuration.DefaultOrDie()
	conf.SetSystemID("10001")
	ds := datasource.New(conf, 10001)

	repo := NewReplyCommentRepository(ds.Orm())

	e, err := repo.GetSubject(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	flog.Info("repo.GetSubject()", "e", e)

	cc := entity.CommentContent{
		Content:    "test comment",
		State:      0,
		CreateBy:   0,
		CreateTime: 0,
	}

	ci, cc, err := repo.CreateComment(context.Background(), 1, 0, 3, 1, 1, cc)

	if err != nil {
		t.Fatal(err)
	}
	flog.Info("repo.CreateComment()", "ci", ci, "cc", cc)

}

func TestListCommentPage(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	conf := configuration.DefaultOrDie()
	conf.SetSystemID("10001")
	ds := datasource.New(conf, 10001)

	ctx := context.Background()
	subject := 1
	parent := 0
	floor := 3
	repo := NewReplyCommentRepository(ds.Orm())

	ids, err := repo.ListCommentsPageIds(context.Background(), uint64(subject), uint64(parent), floor, 10)
	if err != nil {
		t.Fatal(err)
	}
	flog.Info("repo.ListCommentsPageIds(context.Background(), 1, 0, 2, 10)", "ids", ids)

	items, err := repo.ListCommetItem(ctx, ids)
	if err != nil {
		t.Fatal(err)
	}
	flog.Info("repo.ListCommetItem(ctx, ids)", "items", items)

	contents, err := repo.ListCommetContents(ctx, ids)
	if err != nil {
		t.Fatal(err)
	}
	flog.Info("repo.ListCommetContents", "contents", contents)

	n, err := repo.GetCommentLastFloor(ctx, uint64(subject), uint64(parent))
	if err != nil {
		t.Fatal(err)
	}
	flog.Info("repo.GetCommentLastFloor()", "lastfloor", n)
}
