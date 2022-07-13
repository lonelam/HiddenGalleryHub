package models

import "gorm.io/gorm"

type Machine struct {
	gorm.Model
	LatestIp   string
	LatestPort string
	Name       string
	PasswdSum  string
	IsOnline   bool
}
