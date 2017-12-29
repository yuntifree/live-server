package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/dbutil"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/withdraw"
	context "golang.org/x/net/context"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//GetRecords get withdraw records
func (s *Server) GetRecords(ctx context.Context, req *withdraw.GetRequest,
	rsp *withdraw.RecordsResponse) error {
	rows, err := db.Query(`SELECT id, amount, remark, ctime, status FROM 
	withdraw_history WHERE uid = ? ORDER BY id DESC LIMIT ?, ?`,
		req.Uid, req.Seq, req.Num)
	if err == sql.ErrNoRows {
		log.Printf("no more data for uid:%d seq:%d", req.Uid, req.Seq)
		return nil
	} else if err != nil {
		log.Printf("GetRecords query failed:%v", err)
		return err
	}
	defer rows.Close()
	var infos []*withdraw.Record
	for rows.Next() {
		var rec withdraw.Record
		err = rows.Scan(&rec.Id, &rec.Amount, &rec.Remark, &rec.Ctime,
			&rec.Status)
		if err != nil {
			continue
		}
		infos = append(infos, &rec)
	}
	rsp.Infos = infos
	return nil
}

//Apply apply withdraw
func (s *Server) Apply(ctx context.Context, req *withdraw.ApplyRequest,
	rsp *withdraw.ApplyResponse) error {
	var income, apply, withdraw int64
	err := db.QueryRow(`SELECT income, apply, withdraw FROM user_info WHERE uid = ?`,
		req.Uid).Scan(&income, &apply, &withdraw)
	if err != nil {
		log.Printf("Apply query info failed:%d %v", req.Uid, err)
		return err
	}
	if req.Amount+apply+withdraw > income {
		log.Printf("not sufficient remain charge:%d %d", req.Uid, req.Amount)
		return errors.New("not sufficient charge")
	}
	res, err := db.Exec(`INSERT INTO withdraw_history(uid, amount, remark, ctime)
	VALUES (?, ?, ?, NOW())`, req.Uid, req.Amount, req.Remark)
	if err != nil {
		log.Printf("Apply insert withdraw failed:%d %v", req.Uid, err)
		return err
	}
	_, err = db.Exec(`UPDATE user_info SET apply = apply + ? WHERE uid = ?`,
		req.Amount, req.Uid)
	if err != nil {
		log.Printf("Apply update user info failed:%d %v", req.Uid, err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Apply get insert id failed:%d %v", req.Uid, err)
		return err
	}
	rsp.Id = id
	return nil
}

func main() {
	var err error
	db, err = dbutil.NewDB(accounts.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name(accounts.WithdrawService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	withdraw.RegisterWithdrawHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
