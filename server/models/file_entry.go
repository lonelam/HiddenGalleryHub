package models

import "gorm.io/gorm"

type FileEntry struct {
	gorm.Model
	Name              string
	RelativePath      string
	MachineId         uint
	Machine           *Machine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsInvalid         bool
	ParentDirectoryId uint
	ParentDirectory   *Directory
}
