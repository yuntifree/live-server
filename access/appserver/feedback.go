package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/feedback"
)

func feedbackHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "add":
		addFeedback(c)
	case "records":
		getFeedbackRecords(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func addFeedback(c *gin.Context) {
	var info feedback.Info
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	var req feedback.AddRequest
	req.Info = &info
	cl := feedback.NewFeedbackClient(accounts.FeedbackService, client.DefaultClient)
	rsp, err := cl.Add(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", rsp)
	c.JSON(http.StatusOK, gin.H{"errno": 0, "id": rsp.Id})
}

func getFeedbackRecords(c *gin.Context) {
	var req feedback.GetRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	cl := feedback.NewFeedbackClient(accounts.FeedbackService, client.DefaultClient)
	rsp, err := cl.GetRecords(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", rsp)
	c.JSON(http.StatusOK, gin.H{"errno": 0, "infos": rsp.Infos})
}
