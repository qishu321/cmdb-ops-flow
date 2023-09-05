package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/goflow"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"cmdb-ops-flow/utils/ssh"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

var FLW = goflow.New()

func FlowAPI(c *gin.Context) {
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

	// 创建一个 goplow 流程
	flow := goflow.New()

	sort.Slice(list, func(i, j int) bool {
		return list[i].Jobleve < list[j].Jobleve
	})
	// 遍历每个 job
	for _, job := range list {
		job := job // 创建一个本地副本以避免在 goroutine 中的数据竞争
		flow.Add(job.Jobname, []string{}, func(res map[string]interface{}) (interface{}, error) {
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

			return serverResults, nil
		})
	}

	// 执行流程
	results, err := flow.Do()
	//fmt.Println(results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, results, msg.GetErrMsg(200)))
}
