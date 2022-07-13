package routes

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var router = gin.Default()

func RunWithWebsocketUpgrader(host string, pool *connections.Pool, db *gorm.DB) {
	router.GET("/ws/", func(c *gin.Context) {
		pool.AddWsConnection(c.Writer, c.Request, db)
	})
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteString(`<h1>home</h1>`)
	})
	router.GET("/dir/:id", func(c *gin.Context) {
		pageSize, err1 := strconv.Atoi(c.DefaultQuery("page_size", "100"))
		pageIndex, err2 := strconv.Atoi(c.DefaultQuery("page_index", "0"))
		dirID, err3 := strconv.Atoi(c.Param("id"))
		if err1 != nil || err2 != nil || err3 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "page query params is not valid integer",
			})
		}
		subDirArr := make([]models.Directory, 0)
		db = db.Debug()
		db.Where(models.Directory{
			ParentDirectoryId: uint(dirID),
		}).Find(&subDirArr)
		subFileArr := make([]models.FileEntry, 0)
		db.Where(models.FileEntry{
			ParentDirectoryId: uint(dirID),
		}).Order(db.Order("id asc")).Offset(pageSize * pageIndex).Limit(pageSize).Find(&subFileArr)
		c.JSON(http.StatusOK, gin.H{
			"subDirs":  subDirArr,
			"subFiles": subFileArr,
		})
	})
	router.GET("/file/:id", func(c *gin.Context) {
		fileID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "file id is not valid integer",
			})
		}
		var fileEntry models.FileEntry
		db.Preload(clause.Associations).First(&fileEntry, fileID)
		ch, _ := pool.GetFileFromRemote(&fileEntry)
		disconnected := c.Stream(func(w io.Writer) bool {
			buffer, ok := <-ch
			if !ok {
				return false
			}
			w.Write(buffer)
			return true
		})
		if disconnected {
			log.Fatal("file streaming disconnected by client")
		}
	})
	router.Run(host)
}
