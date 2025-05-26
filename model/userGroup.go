package model

import (
	"gorm.io/gorm"
)

type UserGroup struct {
	gorm.Model
	Groupname  string `valid:"matches(^[a-zA-Z][\\w]+$),length(8|16),required" json:"groupName"`
	GroupCount string `valid:"stringlength(8|16),required" json:"groupCount"`
	Group      string `valid:"matches(^1[3-9]{1}\\d{9}$),required" json:"phone"`
}
