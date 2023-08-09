package middleware

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/result"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Token() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("X-Token")
		//fmt.Println("Token value:", token) // 在这里打印 Token 的值
		//fmt.Println(token)
		token_exsits, err := models.TokenInfo(token)
		if err != nil {
			resospnseWithError(401, "非法请求", context)
			return
		}
		fmt.Println(token_exsits)
		if len(token_exsits) != 0 {
			//先做时间判断
			target_time, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Time(token_exsits[0].ExpirationAt).Format("2006-01-02 15:04:05"), time.Local) //需要加上time.Local不然会自动加八小时
			if target_time.Unix() <= time.Now().Unix() {
				//fmt.Println("过期报错")
				//过期报错
				resospnseWithError(401, "timeout", context)
				return
			}
			//token没过期，更新到期时间
			now := time.Unix(time.Now().Unix()+7200, 0).Format("2006-01-02 15:04:05")
			err = models.UpdateTokenTime(token, now)
			context.Set("name", token_exsits[0].Username)
			context.Set("avatar", token_exsits[0].Avatar)
			context.Set("token", token)
		} else {
			fmt.Println("没了")
			resospnseWithError(401, "已退出", context)
			return
		}

		context.Next()
	}
}

func resospnseWithError(code int, message string, c *gin.Context) {
	var res result.Result
	res.Code = code
	res.Msg = message
	c.JSON(200, res) //前端返回也要返回200才能拦截
	c.JSON(http.StatusOK, res)
	c.Abort()
}
