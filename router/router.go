package router

import (
	"cmdb-ops-flow/api"
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
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
	apiv1 := r.Group("api")
	apiv1.Use(middleware.Token())
	{
		apiv1.POST("/user/addUser", api.AddUser)
		apiv1.POST("/user/delUser", api.DelUser)
		apiv1.POST("/user/getUser", api.GetUser)
		apiv1.POST("/user/editUser", api.EditUser)
		apiv1.POST("/user/logout", api.Logout)
		apiv1.POST("/user/info", api.Info)

		apiv1.POST("/cmdb/addCmdb", api.AddCmdb)
		apiv1.POST("/cmdb/getCmdb", api.GetCmdb)
		apiv1.POST("/cmdb/editCmdb", api.EditCmdb)
		apiv1.POST("/cmdb/delCmdb", api.DelCmdb)
		apiv1.POST("/cmdb/GetSearchCmdb", api.GetSearchCmdb)

	}
	adminuser := r.Group("/api/admin/user")
	{
		adminuser.POST("/login", api.Login)

	}

	s.ListenAndServe()
}
