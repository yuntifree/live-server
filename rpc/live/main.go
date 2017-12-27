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

//Stop stop live stream
func (s *Server) Stop(ctx context.Context, req *live.StopRequest,
	rsp *live.StopResponse) error {
	var uid int64
	err := db.QueryRow(`SELECT uid FROM live_history WHERE id = ?`, req.Id).
		Scan(&uid)
	if err != nil {
		log.Printf("Stop query uid failed:%d %v", req.Id, err)
		return err
	}
	if req.Uid != uid {
		log.Printf("Stop uid not matched:%d %d", req.Uid, uid)
		return fmt.Errorf("not matched owner uid")
	}
	_, err = db.Exec(`UPDATE live_history SET ftime = NOW(), status = 2 
	WHERE id = ?`, req.Id)
	if err != nil {
		log.Printf("Stop update history failed:%d %v", req.Id, err)
		return err
	}
	return nil
}

//GetRecords get live records
func (s *Server) GetRecords(ctx context.Context, req *live.GetRequest,
	rsp *live.RecordResponse) error {
	rows, err := db.Query(`SELECT id, title, cover, depict, ctime, ftime, 
	authority, passwd, price, status, replay FROM live_history WHERE uid = ?
	AND id < ? ORDER BY id DESC LIMIT ?`, req.Uid, req.Seq, req.Num)
	if err == sql.ErrNoRows {
		log.Printf("GetRecords no data:%d %d", req.Uid, req.Seq)
		return nil
	} else if err != nil {
		log.Printf("GetRecords query failed:%d %d %v", req.Uid, req.Seq, err)
		return err
	}
	defer rows.Close()
	var infos []*live.Record
	for rows.Next() {
		var info live.Record
		err = rows.Scan(&info.Id, &info.Title, &info.Cover, &info.Depict,
			&info.Ctime, &info.Ftime, &info.Authority, &info.Passwd,
			&info.Price, &info.Status, &info.Replay)
		if err != nil {
			log.Printf("GetRecords scan failed:%v", err)
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
