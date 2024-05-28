package main

import (
	"fmt"
	"github.com/secrity-svc/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/secrity-svc/pkg/logging"
	"github.com/secrity-svc/pkg/setting"
	"github.com/secrity-svc/pkg/util"
	"github.com/secrity-svc/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	util.Setup()
}

// @title Secrity Service API
// @version 1.0
// @description   Secrity Capability
// @license.name MIT
// @license.url https://github.com/secrity-svc/blob/master/LICENSE
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}

}
