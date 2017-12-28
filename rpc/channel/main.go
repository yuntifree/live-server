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
	"github.com/yuntifree/live-server/proto/channel"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//Info return channel info
func (s *Server) Info(ctx context.Context, req *channel.GetRequest,
	rsp *channel.InfoResponse) error {
	var info channel.ChanInfo
	var covers [3]string
	err := db.QueryRow(`SELECT id, title, cover1, cover2, cover3, qrcode, depict,
	chan_intro, live_intro, wxmp, display, dst, extra FROM channel WHERE 
	uid = ?`, req.Uid).Scan(&info.Id, &info.Title, &covers[0], &covers[1],
		&covers[2], &info.Qrcode, &info.Depict, &info.ChanIntro,
		&info.LiveIntro, &info.Wxmp, &info.Display, &info.Dst, &info.Extra)
	if err != nil {
		log.Printf("Info query failed:%d %v", req.Uid, err)
		return err
	}
	info.Covers = covers[:]
	rsp.Info = &info

	return nil
}

func main() {
	var err error
	db, err = dbutil.NewDB(accounts.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(accounts.ChannelService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	channel.RegisterChannelHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
