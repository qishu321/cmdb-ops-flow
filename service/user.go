package service

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"errors"
	"fmt"
	"time"
)

func AddUser(user models.User) (data interface{}, err error) {
	passsword := common.FixMd5(user.Password + conf.Md5Key)
	daoUser := models.User{
		Userid:   common.GenerateRandomUserID(),
		Username: user.Username,
		Password: passsword,
	}

	data, err = models.AddUser(daoUser)
	return data, err
}
func GetUserList(json models.User) (data interface{}, err error) {
	list, err := models.GetUserList(json.ID)
	return list, err
}
func EditUser(user models.User) (data interface{}, err error) {
	passsword := common.FixMd5(user.Password + conf.Md5Key)
	daoUser := models.User{
		Userid:   user.Userid,
		Username: user.Username,
		Password: passsword,
		Avatar:   user.Avatar,
		Role:     user.Role,
	}

	data, err = models.EditUser(daoUser)
	return data, err
}

func Login(json models.User) (data string, err error) {
	//判断有没有这个用户密码
	//没有报错,有了就判断token有没有，没有就创建返回，有就判断时间，时间没过期就返回，过期了就重新生成返回
	var user []models.User
	json.Password = common.FixMd5(json.Password + conf.Md5Key)
	user, err = models.Login(json.Username, json.Password)
	if len(user) > 0 {
		//成功
		if user[0].Token == "" {
			fmt.Println("token是空")
			//token为空，创建更新返回
			now := time.Unix(time.Now().Unix()+7200, 0).Format("2006-01-02 15:04:05")
			token := common.GetRandomString(32)
			err = models.LoginCreateToken(json.Username, json.Password, token, now)
			return token, err
		} else {
			//token不为空 ，判断时间，时间不过期直接返回
			target_time, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Time(user[0].ExpirationAt).Format("2006-01-02 15:04:05"), time.Local) //需要加上time.Local不然会自动加八小时
			if target_time.Unix() >= time.Now().Unix() {
				return user[0].Token, nil
			} else {
				//时间过期
				token := common.GetRandomString(32)
				now := time.Unix(time.Now().Unix()+7200, 0).Format("2006-01-02 15:04:05")
				err = models.LoginCreateToken(json.Username, json.Password, token, now)
				return token, err
			}
		}
	}
	return "登陆失败", errors.New("登陆失败")
}

func Info(token string) (data interface{}, err error) {
	list, err := models.Info(token)
	return list, err
}

func Logout(token string) (err error) {
	err = models.Logout(token)
	return err
}
