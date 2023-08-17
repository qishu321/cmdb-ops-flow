package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"cmdb-ops-flow/utils/ssh"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewCustomAPI(c *gin.Context) {
	var data models.Job
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取 CheckJobgroup 的数据
	list, err := models.CheckJobgroups(data.Jobgroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jobResults []map[string]interface{}

	// 遍历每个 job
	for _, job := range list {
		jobResult := make(map[string]interface{})
		jobResult["jobName"] = job.Jobname

		// 用于存储每个服务器的结果
		serverResults := make(map[string]string)

		// 遍历每个服务器
		for _, cmdb := range job.Cmdbnames {
			// 遍历每个脚本
			for _, script := range job.Scriptnames {
				// 执行命令
				commandOutput, err := ssh.ExecuteRemoteCommand(&script, &cmdb)
				if err != nil {
					serverResults[cmdb.PrivateIP+"_"+script.Name+"_executeCommand"] = err.Error()
				} else {
					serverResults[cmdb.PrivateIP+"_"+script.Name+"_commandOutput"] = "命令执行成功:\n" + commandOutput // 添加命令执行结果
				}

				// 创建文件
				_, err = ssh.CreateRemoteFile(&script, &cmdb)
				if err != nil {
					serverResults[cmdb.PrivateIP+"_"+script.Name+"_createFile"] = err.Error()
				} else {
					serverResults[cmdb.PrivateIP+"_"+script.Name+"_createFile"] = "文件创建成功"
				}
			}
		}

		jobResult["serverResults"] = serverResults
		jobResults = append(jobResults, jobResult)
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, jobResults, msg.GetErrMsg(200)))
}

func CheckJobgroup(c *gin.Context) {
	var data models.Job
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := models.CheckJobgroups(data.Jobgroup)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, list, msg.GetErrMsg(200)))

}

func AddJob(c *gin.Context) {
	var data models.Job
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.CheckJob(data.Jobname)
	if code != msg.SUCCSE {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, "Job name不能重复", msg.GetErrMsg(msg.ERROR)))
		return
	}

	list, err := service.AddJob(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func GetJob(c *gin.Context) {
	//var data models.ScriptManager
	var data struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetJobList(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func EditJob(c *gin.Context) {
	var data models.Job
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.EditJob(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func DelJob(c *gin.Context) {

	var data models.Job
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelJob(data.Jobid)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, msg.SUCCSE, msg.GetErrMsg(code)))

}
