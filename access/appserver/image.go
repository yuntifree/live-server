package main

import (
	"context"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/components/strutil"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/aliyun"
	"github.com/yuntifree/live-server/proto/image"
)

func imageHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "apply":
		imageApply(c)
	case "callback":
		imageFinish(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func imageFinish(c *gin.Context) {
	var req image.FinRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	cl := image.NewImageClient(accounts.ImageService, client.DefaultClient)
	_, err := cl.Finish(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0})
}

type addImageReq struct {
	Uid     int64    `json:"uid"`
	Formats []string `json:"formats"`
}

func imageApply(c *gin.Context) {
	var req addImageReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	var rq image.AddRequest
	rq.Uid = req.Uid
	var names []string
	for i := 0; i < len(req.Formats); i++ {
		name := strutil.GenUUID() + "." + req.Formats[i]
		names = append(names, name)
	}
	rq.Names = names
	cl := image.NewImageClient(accounts.ImageService, client.DefaultClient)
	_, err := cl.Add(context.Background(), &rq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	data, _ := simplejson.NewJson([]byte(`{}`))
	aliyun.FillPolicyResp(data)
	data.Set("errno", 0)
	data.Set("names", names)
	c.JSON(http.StatusOK, data)
}
