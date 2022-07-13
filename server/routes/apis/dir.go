package apis

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddDirApi(router *gin.Engine, pool *connections.Pool, db *gorm.DB) {

	router.GET("api/dir/:id", func(c *gin.Context) {
		pageSize, err1 := strconv.Atoi(c.DefaultQuery("page_size", "100"))
		pageIndex, err2 := strconv.Atoi(c.DefaultQuery("page_index", "0"))
		dirID, err3 := strconv.Atoi(c.Param("id"))
		if err1 != nil || err2 != nil || err3 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "page query params is not valid integer",
			})
		}
		var currentDir models.Directory
		result := db.Preload(clause.Associations).First(&currentDir, dirID)
		if result.RowsAffected <= 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "directory not exist at server",
			})
		}
		subDirArr := make([]models.Directory, 0)
		db.Where(models.Directory{
			ParentDirectoryId: uint(dirID),
		}).Find(&subDirArr)
		subFileArr := make([]models.FileEntry, 0)
		db.Where(models.FileEntry{
			ParentDirectoryId: uint(dirID),
		}).Order(db.Order("id asc")).Offset(pageSize * pageIndex).Limit(pageSize).Find(&subFileArr)
		c.JSON(http.StatusOK, gin.H{
			"Info":     &currentDir,
			"SubDirs":  subDirArr,
			"SubFiles": subFileArr,
		})
	})
}
