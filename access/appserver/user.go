package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/user"
)

const (
	confPath = "/data/live/role.json"
)

func loadRoleConf(path string) *simplejson.Json {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("open %s failed:%v", path, err)
	}
	roleConf, err := simplejson.NewFromReader(file)
	if err != nil {
		log.Fatal("parse json failed:%v", err)
	}
	return roleConf
}

func userHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "login":
		userLogin(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func userLogin(c *gin.Context) {
	var req user.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}

	cl := user.NewUserClient(accounts.UserService, client.DefaultClient)
	rsp, err := cl.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", rsp)
	conf := loadRoleConf(confPath)
	r := fmt.Sprintf("%d", rsp.Role)
	c.JSON(http.StatusOK, gin.H{"errno": 0, "uid": rsp.Uid,
		"token": rsp.Token, "roleconf": conf.Get(r)})
}
