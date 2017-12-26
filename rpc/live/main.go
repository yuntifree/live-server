package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	context "golang.org/x/net/context"

	"github.com/yuntifree/components/dbutil"
	"github.com/yuntifree/components/strutil"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/aliyun"
	live "github.com/yuntifree/live-server/proto/live"
)

var db *sql.DB

//Server server  implement
type Server struct{}

func getUserStream(db *sql.DB, uid int64) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM stream WHERE uid = ?", uid).
		Scan(&name)
	if err != nil {
		//to create new stream
		name = strutil.GenUUID()
		_, err = db.Exec("INSERT INTO stream(name, uid, ctime) VALUES (?, ?, NOW())",
			name, uid)
		if err != nil {
			log.Printf("record stream failed:%v", err)
			return "", err
		}
	}
	return name, nil
}

//Create return a stream name for uid
func (s *Server) Create(ctx context.Context, req *live.CreateRequest,
	rsp *live.CreateResponse) error {
	name, err := getUserStream(db, req.Uid)
	if err != nil {
		log.Printf("Create failed:%d %v", req.Uid, err)
		return err
	}

	push := aliyun.GenPushURL(name, req.Uid)
	_, err = db.Exec(`INSERT INTO live_history(uid, title, cover,
	depict, authority, passwd, price, resolution, push, ctime) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`, req.Uid,
		req.Title, req.Cover, req.Depict, req.Authority, req.Passwd,
		req.Price, req.Resolution, push)
	if err != nil {
		log.Printf("Create insert history failed:%d %v", req.Uid, err)
		return err
	}
	rsp.Push = push
	rsp.Rtmp = aliyun.GenLiveRTMP(name, req.Resolution)
	rsp.Flv = aliyun.GenLiveFLV(name, req.Resolution)
	rsp.Hls = aliyun.GenLiveHLS(name)

	return nil
}

func genAuthStream(name, host, key string, uid int64) string {
	uri := "/live/" + name
	ts := time.Now().Unix() + 3600
	auth := aliyun.GenAuthKey(uri, key, ts, 0, uid)
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
		micro.Name(accounts.LiveService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	live.RegisterLiveHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
