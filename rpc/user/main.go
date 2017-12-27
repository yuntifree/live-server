package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/dbutil"
	"github.com/yuntifree/components/strutil"
	accounts "github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/user"
	context "golang.org/x/net/context"
)

var db *sql.DB

//Server server  implement
type Server struct{}

//Login user login
func (s *Server) Login(ctx context.Context, req *user.LoginRequest,
	rsp *user.LoginResponse) error {
	var uid, role int64
	var pass, salt string
	err := db.QueryRow(`SELECT uid, role, passwd, salt FROM users WHERE name = ?`,
		req.Username).Scan(&uid, &role, &pass, &salt)
	if err != nil {
		log.Printf("Login query failed:%s %v", req.Username, err)
		return err
	}
	epass := strutil.MD5(req.Passwd + salt)
	if epass != pass {
		log.Printf("Login illegal password:%s %s-%s", req.Username, epass, pass)
		return errors.New("illegal password")
	}
	token := strutil.GenSalt()
	_, err = db.Exec("UPDATE users SET token = ? WHERE uid = ?", token, uid)
	if err != nil {
		log.Printf("Login update token failed:%s %v", req.Username, err)
		return err
	}
	rsp.Uid = uid
	rsp.Token = token
	rsp.Role = role
	return nil
}

//Info get user info
func (s *Server) Info(ctx context.Context, req *user.InfoRequest,
	rsp *user.InfoResponse) error {
	err := db.QueryRow(`SELECT name, nickname, headurl, role, ctime, 
	stime, etime, income, expense FROM users u, user_info i WHERE 
	u.uid = i.uid AND u.uid = ?`, req.Uid).Scan(&rsp.Username, &rsp.Nickname,
		&rsp.Headurl, &rsp.Role, &rsp.Ctime, &rsp.Start, &rsp.End,
		&rsp.Income, &rsp.Expense)
	if err != nil {
		log.Printf("Info query failed:%d %v", req.Uid, err)
		return err
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
		micro.Name(accounts.UserService),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()
	user.RegisterUserHandler(service.Server(), new(Server))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
