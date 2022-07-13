package models

import "gorm.io/gorm"

type FileEntry struct {
	gorm.Model
	Name              string
	RelativePath      string
	FileSize          uint
	Thumbnail         string
	ThumbnailHeight   int
	ThumbnailWidth    int
	MachineId         uint
	Machine           *Machine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsInvalid         bool
	ParentDirectoryId uint
	ParentDirectory   *Directory
}
