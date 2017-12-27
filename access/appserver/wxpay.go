package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/live-server/accounts"
	"github.com/yuntifree/live-server/proto/pay"
)

func wxpayHandler(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "scan_callback":
		scanHandle(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
}

func scanHandle(c *gin.Context) {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Read body failed:%v", err)
	}
	log.Printf("body:%s", string(buf))
	var req weixin.ScanReq
	dec := xml.NewDecoder(bytes.NewReader(buf))
	err = dec.Decode(&req)
	log.Printf("req:%+v", req)

	var rq pay.AddRequest
	rq.Appid = req.Appid
	rq.Merid = req.MchID
	rq.Uid, rq.Item = parseProduct(req.ProductID)
	rq.Openid = req.Openid
	rq.Clientip = c.ClientIP()
	cl := pay.NewPayClient(accounts.PayService, client.DefaultClient)
	rsp, err := cl.Add(context.Background(), &rq)
	if err != nil {
		log.Printf("Add failed:%+v %v", req, err)
		c.XML(http.StatusOK, gin.H{"return_code": "FAIL"})
		return
	}
	wx := weixin.WxPay{MerID: accounts.WxMerID,
		MerKey: accounts.WxMerKey}
	var resp weixin.ScanResp
	resp.ReturnCode = "SUCCESS"
	resp.Appid = rsp.Appid
	resp.MchID = rsp.Merid
	resp.NonceStr = rsp.Nonce
	resp.PrepayID = rsp.Prepayid
	resp.ResultCode = "SUCCESS"
	resp.Sign = wx.CalcScanSign(resp)
	out, err := xml.Marshal(resp)
	if err != nil {
		log.Printf("Marshal failed:%v", err)
		c.XML(http.StatusOK, gin.H{"return_code": "FAIL"})
		return
	}
	data := strings.Replace(string(out), "ScanResp", "xml", -1)
	c.String(http.StatusOK, data)
}

func parseProduct(product string) (int64, int64) {
	arr := strings.Split(product, "-")
	if len(arr) != 2 {
		panic("illegal product format " + product)
	}
	uid, err := strconv.Atoi(arr[0])
	if err != nil {
		panic("parseProduct failed:" + err.Error())
	}
	item, err := strconv.Atoi(arr[1])
	if err != nil {
		panic("parseProduct failed:" + err.Error())
	}
	return int64(uid), int64(item)
}
