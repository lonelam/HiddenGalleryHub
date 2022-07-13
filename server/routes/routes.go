package routes

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/routes/apis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router = gin.Default()

func RunWithWebsocketUpgrader(host string, pool *connections.Pool, db *gorm.DB) {
	apis.AddWsApi(router, pool, db)
	apis.AddHomeApi(router, pool, db)
	apis.AddDirApi(router, pool, db)
	apis.AddMachinesApi(router, pool, db)
	apis.AddFileByIdApi(router, pool, db)
	router.StaticFile("/", "../gallery-web-app/build/index.html")
	router.Static("/static", "../gallery-web-app/build")
	router.Run(host)
}
