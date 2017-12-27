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
	"github.com/yuntifree/live-server/proto/image"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//Add add images
func (s *Server) Add(ctx context.Context, req *image.AddRequest,
	rsp *image.AddResponse) error {
	for i := 0; i < len(req.Names); i++ {
		_, err := db.Exec("INSERT IGNORE INTO image(uid, name, ctime) VALUES (?, ?, NOW())",
			req.Uid, req.Names[i])
		if err != nil {
			log.Printf("Add insert failed:%s %v", req.Names[i], err)
		}
	}
	return nil
}

//Finish image upload success
func (s *Server) Finish(ctx context.Context, req *image.FinRequest,
	rsp *image.FinResponse) error {
	_, err := db.Exec(`UPDATE image SET filesize = ?, height = ?,
	width = ?, status = 1, ftime = NOW() WHERE name = ?`, req.Size,
		req.Height, req.Width, req.Filename)
	if err != nil {
		log.Printf("Finish update info failed:%s %v", req.Filename, err)
		return err
	}
	return nil
}

func main() {
	var err error
	db, err = dbutil.NewDB(accounts.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(accounts.ImageService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	image.RegisterImageHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
