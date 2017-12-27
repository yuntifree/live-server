package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	context "golang.org/x/net/context"

	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/dbutil"
	"github.com/yuntifree/components/weixin"
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
		log.Printf("GetRecords no more data for uid:%d seq:%d", req.Uid, req.Seq)
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

//GetRecharges get recharge records
func (s *Server) GetRecharges(ctx context.Context, req *order.GetRequest,
	rsp *order.RechargesResponse) error {
	rows, err := db.Query(`SELECT id, oid, depict, price, ctime, status FROM 
	orders WHERE owner = ? AND type = 1 AND id < ? ORDER BY id DESC LIMIT ?`,
		req.Uid, req.Seq, req.Num)
	if err == sql.ErrNoRows {
		log.Printf("GetRecharges no more data for uid:%d seq:%d", req.Uid, req.Seq)
		return nil
	} else if err != nil {
		log.Printf("GetRecharges query failed:%v", err)
		return err
	}
	defer rows.Close()
	var infos []*order.Recharge
	for rows.Next() {
		var info order.Recharge
		err = rows.Scan(&info.Id, &info.Oid, &info.Depict, &info.Price,
			&info.Ctime, &info.Status)
		if err != nil {
			continue
		}
		infos = append(infos, &info)
	}
	rsp.Infos = infos
	return nil
}

//GetItems get pay items
func (s *Server) GetItems(ctx context.Context, req *order.GetRequest,
	rsp *order.ItemsResponse) error {
	rows, err := db.Query(`SELECT id, price, img FROM pay_items
	WHERE deleted = 0 AND online = 1`)
	if err != nil {
		log.Printf("GetItems query failed:%v", err)
		return err
	}
	defer rows.Close()
	var infos []*order.Item
	wx := weixin.WxPay{MerID: accounts.WxMerID,
		MerKey: accounts.WxMerKey, Appid: accounts.DgWxAppid}
	for rows.Next() {
		var item order.Item
		err = rows.Scan(&item.Id, &item.Price,
			&item.Img)
		if err != nil {
			continue
		}
		product := fmt.Sprintf("%d-%d", req.Uid, item.Id)
		item.Qrcode = wx.GenQRCode(product)
		infos = append(infos, &item)
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
