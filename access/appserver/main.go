package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuntifree/components/weixin"
)

func scanHandler(c *gin.Context) {
	log.Printf("url:%+v", c.Request)
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Read body failed:%v", err)
	}
	log.Printf("body:%s", string(buf))
	var req weixin.ScanReq
	dec := xml.NewDecoder(bytes.NewReader(buf))
	err = dec.Decode(&req)
	log.Printf("req:%+v", req)
}

func main() {
	router := gin.Default()
	router.POST("/user/:action", userHandler)
	router.POST("/withdraw/:action", withdrawHandler)
	router.POST("/order/:action", orderHandler)
	router.POST("/live/:action", liveHandler)
	router.POST("/image/:action", userHandler)
	router.POST("/wxpay/scan_callback", scanHandler)

	srv := &http.Server{
		Addr:    ":9898",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen:%s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
