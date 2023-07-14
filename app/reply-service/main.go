package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/jinvei/microservice/app/reply-service/domain/repository"
	"github.com/jinvei/microservice/app/reply-service/domain/service"
	"github.com/jinvei/microservice/base/api/codes"
	"github.com/jinvei/microservice/base/api/proto/v1/app"
	"github.com/jinvei/microservice/base/framework/configuration"
	"github.com/jinvei/microservice/base/framework/datasource"
	"github.com/jinvei/microservice/base/framework/rpc"
	"google.golang.org/grpc"
)

func main() {
	conf := configuration.DefaultOrDie()
	systemID := codes.ReplySvcSystemID
	conf.SetSystemID(strconv.Itoa(systemID))

	db := datasource.New(conf, systemID)
	repo := repository.NewReplyCommentRepository(db.Gorm(), uint64(systemID))
	replysvc := service.NewReplyCommentService(repo, conf)

	var waitSrv sync.WaitGroup
	waitSrv.Add(1)
	go func() {
		defer waitSrv.Done()

		err := rpc.Serve(conf, systemID, func(srv *grpc.Server) {
			app.RegisterReplyCommentServiceServer(srv, replysvc)
		})

		if err != nil {
			log.Println(err)
		}
	}()

	waitSrv.Wait()
}
