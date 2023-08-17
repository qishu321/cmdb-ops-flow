package models

import (
	"time"
)

type Jobhistory struct {
	ID             int        `json:"id" db:"id" form:"id"`
	Jobhistoryid   int64      `gorm:"type:bigint;not null" json:"Jobhistoryid" validate:"required"`
	Jobhistoryname string     `json:"jobhistoryname" db:"jobhistoryname" form:"jobhistoryname"`
	Jobgroup       string     `json:"jobgroup" db:"jobgroup" form:"jobgroup"`
	Jobgroups      []JobGroup `gorm:"FOREIGNKEY:Jobgroupname;ASSOCIATION_FOREIGNKEY:Jobgroup"`
	Jobnames       string     `json:"jobname" db:"jobname" form:"jobname"`
	Jobname        []Job      `gorm:"FOREIGNKEY:Jobname;ASSOCIATION_FOREIGNKEY:Jobnames"`
	start_time     time.Time
	end_time       time.Time
	Status         string    `json:"status" db:"status" form:"status"`
	Jobhistorylog  string    `json:"jobhistorylog" db:"jobhistorylog" form:"jobhistorylog"`
	created_at     time.Time `gorm:"default:CURRENT_TIMESTAMP",json:"created_at"`
	Label          string    `json:"label" db:"label" form:"label"`
}
