package models

import "gorm.io/gorm"

type Directory struct {
	gorm.Model
	Name              string
	RelativePath      string   `gorm:"uniqueIndex:dir_combined_unique;not null;"`
	MachineId         uint     `gorm:"uniqueIndex:dir_combined_unique;not null;"`
	Machine           *Machine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsRootDirectory   bool
	IsInvalid         bool
	ParentDirectoryId uint
	ParentDirectory   *Directory
}
