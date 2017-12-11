package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/stream"
)

func streamHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "get":
		getStream(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func getStream(c *gin.Context) {
	var req stream.CreateRequest
	u := c.Query("uid")
	uid, _ := strconv.Atoi(u)
	req.Uid = int64(uid)
	cl := stream.NewStreamClient(accounts.StreamService, client.DefaultClient)
	rsp, err := cl.Create(context.Background(), &req)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "url": rsp.Url,
		"stream": rsp.Stream})
}
