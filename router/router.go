package router

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/middleware"
	"net/http"
	"github.com/gin-gonic/gin"

)

func InitRouter() {
	gin.SetMode(conf.AppMode)
	r := gin.Default()
	//fmt.Println(utils.HttpPort)
	r.Use(middleware.Cors())

	s := &http.Server{
		Addr:           conf.HttpPort,
		Handler:        r,
		ReadTimeout:    conf.ReadTimeout,
		WriteTimeout:   conf.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "PONG")
	})
	s.ListenAndServe()
}