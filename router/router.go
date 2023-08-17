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

		apiv1.POST("/cmdb/ssh/command", api.SshCommand)
		apiv1.POST("/cmdb/ssh/createFile", api.SshcreateFile)

		apiv1.GET("/cmdb/ssh/webssh", api.VisitorWebsocketServer)

		apiv1.POST("/script/addScript", api.AddScript)
		apiv1.POST("/script/getScript", api.GetScript)
		apiv1.POST("/script/editScript", api.EditScript)
		apiv1.POST("/script/delScript", api.DelScript)

		apiv1.POST("/job/addJob", api.AddJob)
		apiv1.POST("/job/getJob", api.GetJob)
		apiv1.POST("/job/editJob", api.EditJob)
		apiv1.POST("/job/delJob", api.DelJob)
		apiv1.POST("/job/CheckJobgroup", api.CheckJobgroup)
		apiv1.POST("/job/NewCustomAPI", api.NewCustomAPI)

		apiv1.POST("/job/Group/addJobGroup", api.AddJobGroup)
		apiv1.POST("/job/Group/GetJobGroup", api.GetJobGroup)
		apiv1.POST("/job/Group/EditJobGroup", api.EditJobGroup)
		apiv1.POST("/job/Group/DelJobGroup", api.DelJobGroup)
		apiv1.POST("/job/Group/GetSearchJobGroup", api.GetSearchJobGroup)

	}
	adminuser := r.Group("/api/admin/user")
	{
		adminuser.POST("/login", api.Login)
		adminuser.GET("/ssh/webssh", api.VisitorWebsocketServer)

	}

	s.ListenAndServe()
}
