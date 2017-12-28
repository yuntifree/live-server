package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/channel"
)

func channelHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "info":
		getChannelInfo(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func getChannelInfo(c *gin.Context) {
	var req channel.GetRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	cl := channel.NewChannelClient(accounts.ChannelService, client.DefaultClient)
	rsp, err := cl.Info(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "info": rsp.Info})
}
