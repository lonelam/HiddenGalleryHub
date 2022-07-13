package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Machine{})
	db.AutoMigrate(&Directory{})
	db.AutoMigrate(&FileEntry{})
}
