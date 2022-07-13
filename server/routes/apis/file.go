package apis

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var IMAGE_EXTS2CONTENT_TYPE = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".jfif": "image/jpeg",
}

func AddFileInfoByIdApi(router *gin.Engine, pool *connections.Pool, db *gorm.DB) {
	router.GET("api/file_info/:id", func(c *gin.Context) {
		fileID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "file id is not valid integer",
			})
			return
		}
		var fileEntry models.FileEntry
		result := db.Preload(clause.Associations).First(&fileEntry, fileID)
		if result.RowsAffected <= 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "file not found at server",
			})
			return
		}
		c.JSON(http.StatusOK, &fileEntry)
	})
}

func AddFileByIdApi(router *gin.Engine, pool *connections.Pool, db *gorm.DB) {

	router.GET("api/file/:id/*", func(c *gin.Context) {
		fileID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "file id is not valid integer",
			})
		}
		var fileEntry models.FileEntry
		db.Preload(clause.Associations).First(&fileEntry, fileID)
		ch, _ := pool.GetFileFromRemote(&fileEntry)
		contentLength := fileEntry.FileSize
		// contentType := "application/octet-stream"
		contentDisposition := fmt.Sprintf(`attachment; filename="%s"`, fileEntry.Name)
		if _, hasKey := IMAGE_EXTS2CONTENT_TYPE[strings.ToLower(path.Ext(fileEntry.Name))]; hasKey {
			// contentType = IMAGE_EXTS2CONTENT_TYPE[strings.ToLower(path.Ext(fileEntry.Name))]
		} else {
			c.Header("Content-Disposition", contentDisposition)
		}

		// c.Header("Content-Type", contentType)
		c.Header("Content-Length", strconv.FormatUint(uint64(contentLength), 10))
		done := false
		respBuffer := make([]byte, 0)
		mu := sync.Mutex{}
		go func() {
			for {
				buffer, ok := <-ch
				if !ok {
					done = true
					return
				}
				mu.Lock()
				respBuffer = append(respBuffer, buffer...)
				mu.Unlock()
			}
		}()
		for {
			if len(respBuffer) > 0 {
				mu.Lock()
				writtenBuffer := respBuffer
				respBuffer = make([]byte, 0)
				mu.Unlock()
				c.Writer.Write(writtenBuffer)
			}

			if done {
				break
			}
		}
	})
}
