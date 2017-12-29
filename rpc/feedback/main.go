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
	res, err := db.Exec(`INSERT INTO feedback(uid, title, content, img, ctime)
	VALUES (?, ?, ?, ?, NOW())`, req.Info.Uid, req.Info.Title, req.Info.Content,
		req.Info.Img)
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

//GetRecords get feedback records
func (s *Server) GetRecords(ctx context.Context, req *feedback.GetRequest,
	rsp *feedback.RecordsResponse) error {
	rows, err := db.Query(`SELECT id, title, content, img, status FROM feedback 
	WHERE uid = ? ORDER BY id DESC LIMIT ?, ?`, req.Uid,
		req.Seq, req.Num)
	if err != nil {
		log.Printf("GetRecords query failed:%v", err)
		return err
	}
	defer rows.Close()
	var infos []*feedback.Info
	for rows.Next() {
		var info feedback.Info
		err = rows.Scan(&info.Id, &info.Title, &info.Content,
			&info.Img, &info.Status)
		if err != nil {
			continue
		}
		infos = append(infos, &info)
	}
	rsp.Infos = infos
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
