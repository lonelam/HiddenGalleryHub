package apis

import (
	"HiddenGalleryHub/server/connections"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddHomeApi(router *gin.Engine, pool *connections.Pool, db *gorm.DB) {
	router.GET("api/", func(c *gin.Context) {
		c.Writer.WriteString(`<h1>home</h1>`)
	})
}
