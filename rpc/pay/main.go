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
	VALUES (1, ?, ?, ?, ?, NOW())`, oid, req.Uid, req.Item, price)
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

//Fin pay finished
func (s *Server) Fin(ctx context.Context, req *pay.FinRequest,
	rsp *pay.FinResponse) error {
	log.Printf("Fin request:%+v", req)
	var id, item, uid, price, status, typ int64
	var prepayid string
	err := db.QueryRow(`SELECT id, item, uid, price, type, status, prepayid
	FROM orders WHERE oid = ?`, req.Oid).
		Scan(&id, &item, &uid, &price, &typ, &status, &prepayid)
	if err != nil {
		log.Printf("Fin query order info failed:%s %v", req.Oid, err)
		return err
	}
	log.Printf("id:%d item:%d price:%d status:%d prepayid:%s", id,
		item, price, status, prepayid)
	if status == 1 {
		log.Printf("Fin has duplicated oid:%s", req.Oid)
		return nil
	}
	if price > req.Fee {
		log.Printf("Fin illegal fee, oid:%s %d-%d", req.Oid, price, req.Fee)
		return fmt.Errorf("illegal feed oid:%s %d-%d", req.Oid, price, req.Fee)
	}
	_, err = db.Exec(`UPDATE orders SET status = 1, fee = ?, ptime = NOW() 
	WHERE id = ?`, req.Fee, id)
	if err != nil {
		log.Printf("Fin update order info failed, oid:%s fee:%d %v", req.Oid,
			req.Fee, err)
		return fmt.Errorf("update order info failed, oid:%s fee:%d %v", req.Oid,
			req.Fee, err)
	}

	log.Printf("after update orders status:%s", req.Oid)
	_, err = db.Exec(`UPDATE user_info SET recharge = recharge + ? WHERE uid = ?`,
		req.Fee, uid)
	if err != nil {
		log.Printf("update user recharge failed:%d %v", uid, err)
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
