package models

import (
	"cmdb-ops-flow/utils/msg"
	"time"
)

type User struct {
	ID       int    `gorm:"primary_key"`
	Userid   int64  `gorm:"type:bigint;not null" json:"userid" validate:"required"`
	Username string `json:"username" db:"username" form:"username" ` // Make username required
	Password string `json:"password" db:"password" form:"password" ` // Make password required

	Avatar       string `json:"avatar" db:"avatar" form:"avatar"`
	Token        string
	Role         int       `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
	ExpirationAt time.Time `gorm:"default:CURRENT_TIMESTAMP",json:"created_at"`
}

func AddUser(user User) (interface{}, error) {

	err := db.Create(&user).Error
	return user, err
}
func EditUser(user User) (interface{}, error) {
	err := db.Select("id").Where("username = ? AND userid = ?", user.Username, user.Userid).First(&user).Error
	if err != nil {
		return msg.ERROR_USERNAME_USED, err
	}
	err = db.Model(&user).Where("username = ? AND userid = ?", user.Username, user.Userid).Updates(user).Error
	if err != nil {
		return user, err
	}

	return user, err
}

func DelUser(name int64) (code int) {
	var user User
	db.Select("id").Where("userid = ?", name).First(&user)
	if user.ID > 0 {
		err = db.Where("userid = ?", name).Delete(&user).Error
		if err != nil {
			return msg.ERROR_USER_NOT_EXIST
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR
	}

}
func GetUserList(id int) ([]User, error) {
	var list []User
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckUser(name string) (code int) {
	var user User
	db.Select("id").Where("username = ?", name).First(&user)
	if user.ID > 0 {
		return msg.ERROR_USERNAME_USED
	}
	return msg.SUCCSE
}
func Login(name string, password string) (list []User, err error) {
	var user []User
	db.Debug().Where("username = ? and password = ?", name, password).First(&user)

	return user, nil
}
func LoginCreateToken(name string, password string, token string, expiration_at string) (err error) {
	return db.Debug().Table("user").Where("username = ? and password = ?", name, password).Updates(map[string]interface{}{"token": token, "expiration_at": expiration_at}).Error
}
func TokenInfo(token string) (list []User, err error) {
	var user []User
	db.Debug().Where("token = ? ", token).First(&user)
	return user, nil
}
func UpdateTokenTime(token string, expiration_at string) (err error) {
	return db.Debug().Table("user").Where("token = ? ", token).Updates(map[string]interface{}{"expiration_at": expiration_at}).Error
}

func Info(token string) (data interface{}, err error) {
	type Result struct {
		Username string
		Avatar   string
	}
	var result Result
	db.Debug().Table("user").Select("username, avatar").Where("token = ? ", token).Scan(&result)
	return result, nil
}
func Logout(token string) (err error) {
	return db.Debug().Table("user").Where("token = ? ", token).Updates(map[string]interface{}{"token": ""}).Error
}

//// CheckUser 查询密码是否存在
//func Check_Login(user User) (code int, err error) {
//	err = db.Select("id").Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return 500, err
//		}
//		return 500, err
//	}
//	return msg.SUCCSE, err
//}
