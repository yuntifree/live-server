package main

import (
	"database/sql"
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
