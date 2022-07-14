package models

import "gorm.io/gorm"

type FileEntry struct {
	gorm.Model
	Name              string
	RelativePath      string `gorm:"uniqueIndex:file_combined_unique;not null;"`
	FileSize          uint
	Thumbnail         string
	ThumbnailHeight   int
	ThumbnailWidth    int
	MachineId         uint     `gorm:"uniqueIndex:file_combined_unique;not null;"`
	Machine           *Machine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsInvalid         bool
	ParentDirectoryId uint
	ParentDirectory   *Directory
}
