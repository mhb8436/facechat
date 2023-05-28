package api

import (
	"context"
	"facechat/api/router"
	"facechat/api/rpc"
	"facechat/config"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Api struct {
}

func New() *Api {
	return &Api{}
}

func (a *Api) Run() {
	// init rpc client
	rpc.InitMessageLogicRpcClient()

	r := router.Register()
	runMode := config.GetGinRunMode()
	logrus.Info("server start, now run mode is", runMode)
	gin.SetMode(runMode)
	apiConfig := config.Conf.Api
	port := apiConfig.ApiBase.ListenPort
	flag.Parse()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("start listen : %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Infof("Shutdown Server ....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")
	os.Exit(0)
}
