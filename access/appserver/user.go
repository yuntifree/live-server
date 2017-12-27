package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/order"
	"github.com/yuntifree/live-server/proto/user"
)

func userHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "login":
		userLogin(c)
	case "summary":
		getUserSummary(c)
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
	c.JSON(http.StatusOK, gin.H{"errno": 0, "uid": rsp.Uid,
		"token": rsp.Token, "role": rsp.Role})
}

func getUserSummary(c *gin.Context) {
	var req commReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	var rq user.InfoRequest
	rq.Uid = req.Uid
	ul := user.NewUserClient(accounts.UserService, client.DefaultClient)
	ursp, err := ul.Info(context.Background(), &rq)
	if err != nil {
		log.Printf("getUserSummary get user info failed:%v", err)
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", ursp)

	var orq order.GetRequest
	orq.Uid = req.Uid
	ol := order.NewOrderClient(accounts.OrderService, client.DefaultClient)
	orsp, err := ol.GetItems(context.Background(), &orq)
	if err != nil {
		log.Printf("getUserSummary get order items failed:%v", err)
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", orsp)

	c.JSON(http.StatusOK, gin.H{"errno": 0, "userinfo": ursp,
		"items": orsp.Infos})
}
