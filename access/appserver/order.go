package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/order"
)

func orderHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "records":
		getOrderRecords(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func getOrderRecords(c *gin.Context) {
	var req order.GetRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 102, "desc": "illegal param"})
		return
	}
	cl := order.NewOrderClient(accounts.OrderService, client.DefaultClient)
	rsp, err := cl.GetRecords(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	log.Printf("rsp:%+v", rsp)
	c.JSON(http.StatusOK, gin.H{"errno": 0, "infos": rsp.Infos})
}
