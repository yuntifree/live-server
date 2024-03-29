package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user/:action", userHandler)
	router.POST("/withdraw/:action", withdrawHandler)
	router.POST("/order/:action", orderHandler)
	router.POST("/live/:action", liveHandler)
	router.POST("/image/:action", imageHandler)
	router.POST("/wxpay/:action", wxpayHandler)
	router.POST("/channel/:action", channelHandler)
	router.POST("/feedback/:action", feedbackHandler)
	router.GET("/live/:action", notifyHandler)

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
