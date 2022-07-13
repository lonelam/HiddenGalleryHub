package apis

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddMachinesApi(router *gin.Engine, pool *connections.Pool, db *gorm.DB) {

	router.GET("api/machines", func(c *gin.Context) {
		var rootDirs []models.Directory
		db.Preload(clause.Associations).Model(&models.Directory{}).Where(&models.Directory{
			IsRootDirectory: true,
		}).Find(&rootDirs)

		c.JSON(http.StatusOK, gin.H{
			"machines": rootDirs,
		})
	})
}
