package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/live"
)

func liveHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "create":
		createLive(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func createLive(c *gin.Context) {
	var req live.CreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	cl := live.NewLiveClient(accounts.LiveService, client.DefaultClient)
	rsp, err := cl.Create(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", rsp)
	c.JSON(http.StatusOK, gin.H{"errno": 0, "push": rsp.Push,
		"rtmp": rsp.Rtmp, "flv": rsp.Flv, "hls": rsp.Hls})
}
