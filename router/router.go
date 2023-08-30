package router

import (
	"cmdb-ops-flow/api"
	"cmdb-ops-flow/api/apis_k8s"
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	//r.Use(middleware.Jaeger())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "PONG")
	})
	r.GET("/metrics", func(handler http.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			handler.ServeHTTP(c.Writer, c.Request)
		}
	}(promhttp.Handler()))

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

		apiv1.POST("/etcd/etcdGetall", api.EtcdGetall)
		apiv1.POST("/etcd/Etcdrestore", api.Etcdrestore)

		apiv1.POST("/etcd/getEtcd", api.GetEtcd)
		apiv1.POST("/etcd/addEtcd", api.AddEtcd)

		apiv1.POST("/etcd/editEtcd", api.EditEtcd)
		apiv1.POST("/etcd/delEtcd", api.DelEtcd)

		apiv1.POST("/etcd/getEtcdbak", api.GetEtcdbak)
		apiv1.POST("/etcd/delEtcdbak", api.DelEtcdbak)
		apiv1.POST("/etcd/addEtcdbak", api.AddEtcdbak)

		apiv1.POST("/job/addJob", api.AddJob)
		apiv1.POST("/job/getJob", api.GetJob)
		apiv1.POST("/job/editJob", api.EditJob)
		apiv1.POST("/job/delJob", api.DelJob)
		apiv1.POST("/job/CheckJobgroup", api.CheckJobgroup)
		apiv1.POST("/job/NewCustomAPI", api.NewCustomAPI)

		//apiv1.POST("/kube/config/addconfig", api.AddKubeConfig)
		//apiv1.POST("/kube/config/getconfig", api.GetKubeConfig)
		//apiv1.POST("/kube/config/editconfig", api.EditKubeConfig)
		//apiv1.POST("/kube/config/delconfig", api.DelKubeConfig)

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
		adminuser.GET("/kube/pods/SshPod", apis_k8s.SshPod)
		adminuser.GET("/kube/pods/getPodLogs", apis_k8s.GetPodLogs)

	}

	api_k8s := r.Group("/api/k8s/")
	api_k8s.Use(middleware.Token())
	{
		api_k8s.POST("/kube/config/addconfig", api.AddKubeConfig)
		api_k8s.POST("/kube/config/getconfig", api.GetKubeConfig)
		api_k8s.POST("/kube/config/editconfig", api.EditKubeConfig)
		api_k8s.POST("/kube/config/delconfig", api.DelKubeConfig)

		api_k8s.POST("/kube/pods/getallPods", apis_k8s.GetAllPods)
		api_k8s.POST("/kube/pods/getPods", apis_k8s.GetPods)
		api_k8s.POST("/kube/pods/getPodLogs", apis_k8s.GetPodLogs)
		api_k8s.POST("/kube/pods/GetPodsYaml", apis_k8s.GetPodsYaml)

		api_k8s.GET("/kube/pods/SshPod", apis_k8s.SshPod)

		api_k8s.POST("/kube/nodes/getVersion", apis_k8s.GetVersion)
		api_k8s.POST("/kube/nodes/getNodeMetrics", apis_k8s.GetNodeMetrics)
		api_k8s.POST("/kube/nodes/getallNodes", apis_k8s.GetAllNodes)

		api_k8s.POST("/kube/deploy/getDeployment", apis_k8s.GetDeployment)

		api_k8s.POST("/kube/ns/getallNamespace", apis_k8s.GetallNamespace)
		api_k8s.POST("/kube/ns/addNamespace", apis_k8s.AddNamespace)
		api_k8s.POST("/kube/ns/editNamespace", apis_k8s.EditNamespace)
		api_k8s.POST("/kube/ns/delNamespace", apis_k8s.DelNamespace)

		api_k8s.POST("/kube/svc/getallSvc", apis_k8s.GetallSvc)
		api_k8s.POST("/kube/svc/delSvc", apis_k8s.DelSvc)
		api_k8s.POST("/kube/svc/addSvc", apis_k8s.AddSvc)
		api_k8s.POST("/kube/svc/editSvc", apis_k8s.EditSvc)

	}

	s.ListenAndServe()
}
