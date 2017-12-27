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
	"github.com/yuntifree/components/strutil"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/pay"
)

const (
	tradeType = "NATIVE"
	succCode  = "SUCCESS"
	signType  = "MD5"
	cbURL     = "http://video.yunxingzh.com/wxpay/callback"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//Add add images
func (s *Server) Add(ctx context.Context, req *pay.AddRequest,
	rsp *pay.AddResponse) error {
	var price int64
	err := db.QueryRow(`SELECT price FROM pay_items WHERE id = ?`, req.Item).
		Scan(&price)
	if err != nil {
		return err
	}
	oid := weixin.GenOrderID(req.Uid)
	id, err := recordOrderInfo(db, oid, price, req)
	if err != nil {
		log.Printf("Add recordOrderInfo failed:%+v %v", req, err)
		return err
	}

	var rq weixin.UnifyOrderReq
	rq.Appid = req.Appid
	rq.Body = "充值"
	rq.MchID = req.Merid
	rq.NonceStr = strutil.GenNonceStr()
	rq.Openid = req.Openid
	rq.TradeType = tradeType
	rq.SpbillCreateIP = req.Clientip
	rq.TotalFee = price
	rq.OutTradeNO = oid
	rq.NotifyURL = cbURL

	wx := weixin.WxPay{MerID: accounts.WxMerID,
		MerKey: accounts.WxMerKey}
	resp, err := wx.UnifyPayRequest(rq)
	if err != nil {
		log.Printf("Add UnifyPayRequest failed:%v", err)
		return err
	}
	log.Printf("resp:%+v", resp)
	if resp.ReturnCode != succCode || resp.ResultCode != succCode {
		log.Printf("WxPay UnifyPayRequest failed msg:%s", resp.ReturnMsg)
		return fmt.Errorf("pay failed:%s", resp.ReturnMsg)
	}
	recordPrepayid(db, id, resp.PrepayID)
	rsp.Appid = req.Appid
	rsp.Merid = req.Merid
	rsp.Nonce = resp.NonceStr
	rsp.Prepayid = resp.PrepayID
	return nil
}

func recordOrderInfo(db *sql.DB, oid string, price int64, req *pay.AddRequest) (int64, error) {
	res, err := db.Exec(`INSERT INTO orders(type, oid, uid, item, price, ctime)
	VALUES (0, ?, ?, ?, ?, NOW())`, oid, req.Uid, req.Item, price)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func recordPrepayid(db *sql.DB, id int64, prepayid string) {
	_, err := db.Exec("UPDATE orders SET prepayid = ? WHERE id = ?",
		prepayid, id)
	if err != nil {
		log.Printf("recordPrepayid failed:%d %v", id, err)
	}
}

func main() {
	var err error
	db, err = dbutil.NewDB(accounts.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(accounts.PayService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	pay.RegisterPayHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
