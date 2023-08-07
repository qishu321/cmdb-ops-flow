package models

type Cmdb struct {
	ID         int    `json:"id" db:"id" form:"id"`
	Cmdbid     int64  `gorm:"type:bigint;not null" json:"cmdbid" validate:"required"`
	Cmdbname   string `json:"cmdbname" db:"cmdbname" form:"cmdbname"`
	PublicIP   string `json:"public_ip" db:"public_ip" form:"public_ip"`
	PrivateIP  string `json:"private_ip" db:"private_ip" form:"private_ip"`
	Username   string `json:"username" db:"username" form:"username"`
	Password   string `json:"password" db:"password" form:"password"`
	PrivateKey string `json:"private_key" db:"private_key" form:"private_key"`
	SSHPort    int    `json:"ssh_port" db:"ssh_port" form:"ssh_port"`
	Label		string	`json:"label" db:"label" form:"label"`
}