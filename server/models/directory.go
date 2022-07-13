package models

import "gorm.io/gorm"

type Directory struct {
	gorm.Model
	Name              string
	RelativePath      string
	MachineId         uint
	Machine           *Machine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsRootDirectory   bool
	IsInvalid         bool
	ParentDirectoryId uint
	ParentDirectory   *Directory
}
