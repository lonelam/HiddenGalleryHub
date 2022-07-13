package apis

import (
	"HiddenGalleryHub/server/connections"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddWsApi(router *gin.Engine, pool *connections.Pool, db *gorm.DB) {

	router.GET("ws/", func(c *gin.Context) {
		pool.AddWsConnection(c.Writer, c.Request, db)
	})
}
