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
	"github.com/yuntifree/live-server/proto/order"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//GetRecords get withdraw records
func (s *Server) GetRecords(ctx context.Context, req *order.GetRequest,
	rsp *order.RecordsResponse) error {
	rows, err := db.Query(`SELECT o.id, o.hid, u.headurl, u.nickname, o.depict, o.price,
	o.uid FROM orders o, users u WHERE o.uid = u.uid AND o.owner = ? AND type = 0 
	AND status = 1 ORDER BY id DESC LIMIT ?, ?`, req.Uid, req.Seq, req.Num)
	if err == sql.ErrNoRows {
		log.Printf("no more data for uid:%d seq:%d", req.Uid, req.Seq)
		return nil
	} else if err != nil {
		log.Printf("GetRecords query failed:%v", err)
		return err
	}
	defer rows.Close()
	var infos []*order.Record
	for rows.Next() {
		var rec order.Record
		err = rows.Scan(&rec.Id, &rec.Hid, &rec.Headurl, &rec.Nickname,
			&rec.Depict, &rec.Price, &rec.Uid)
		if err != nil {
			continue
		}
		infos = append(infos, &rec)
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
		micro.Name(accounts.OrderService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	order.RegisterOrderHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
