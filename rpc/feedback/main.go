package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	context "golang.org/x/net/context"

	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/dbutil"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/feedback"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//Add add feedback
func (s *Server) Add(ctx context.Context, req *feedback.AddRequest,
	rsp *feedback.AddResponse) error {
	res, err := db.Exec(`INSERT INTO feedback(uid, title, content, ctime)
	VALUES (?, ?, ?, NOW())`, req.Info.Uid, req.Info.Title, req.Info.Content)
	if err != nil {
		log.Printf("Add insert failed:%+v %v", req, err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Add get insert id failed:%v", err)
		return err
	}
	rsp.Id = id
	return nil
}

func main() {
	var err error
	db, err = dbutil.NewDB(accounts.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(accounts.FeedbackService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	feedback.RegisterFeedbackHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
