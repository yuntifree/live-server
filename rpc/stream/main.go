package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/dbutil"
	"github.com/yuntifree/components/strutil"
	accounts "github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/live"
	"github.com/yuntifree/live-server/proto/stream"
	context "golang.org/x/net/context"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//Create return a stream name for uid
func (s *Server) Create(ctx context.Context, req *stream.CreateRequest,
	rsp *stream.CreateResponse) error {
	var name string
	err := db.QueryRow("SELECT name FROM stream WHERE uid = ?", req.Uid).
		Scan(&name)
	if err != nil {
		//to create new stream
		name = strutil.GenUUID()
		_, err = db.Exec("INSERT INTO stream(name, uid, ctime) VALUES (?, ?, NOW())",
			name, req.Uid)
		if err != nil {
			log.Printf("record stream failed:%v", err)
			return err
		}
	}

	rsp.Url = accounts.LiveURL
	rsp.Stream = genAuthStream(name, accounts.LiveHost, accounts.AuthKey,
		req.Uid)
	return nil
}

func genAuthStream(name, host, key string, uid int64) string {
	uri := "/live/" + name
	ts := time.Now().Unix() + 3600
	auth := live.GenAuthKey(uri, key, ts, 0, uid)
	return fmt.Sprintf("%s?vhost=%s&auth_key=%d-0-%d-%s",
		name, host, ts, uid, auth)
}

func main() {
	var err error
	db, err = dbutil.NewDB(accounts.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(accounts.StreamService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	stream.RegisterStreamHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
