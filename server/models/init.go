package models

import "gorm.io/gorm"

func Init(db *gorm.DB) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Machine{}).UpdateColumn("is_online", 0)

}
