package processor

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/elliotchance/pie/v2"
)

func (c *WsClientConnection) onRequestDirectorySendAllDirectories(message []byte) {
	requestDirectory, err := messages.ReadRequestDirectoryMessage(message)
	if err != nil {
		log.Fatal(err)
		return
	}
	startPath := filepath.Join(c.rootDir, requestDirectory.RelativePath)
	startRelPath, _ := filepath.Rel(c.rootDir, startPath)

	directoryArr := make([]messages.DirectoryEntry, 1)
	directoryArr[0] = messages.DirectoryEntry{
		Name:               path.Base(startPath),
		RelativePath:       startRelPath,
		ParentRelativePath: path.Dir(startRelPath),
	}
	fileArr := make([]messages.FileEntry, 0)
	// log.Printf("root:%s relative: %s startpath: %s\n", c.rootDir, requestDirectory.RelativePath, startPath)
	collectDirectoryStructure(c.rootDir, startPath, &directoryArr, &fileArr, false)
	dirStruct := messages.DirectoryStructureMessage{
		DirectoryEntries: directoryArr,
		FileEntries:      fileArr,
	}
	dirStructMsg, _ := json.Marshal(dirStruct)
	ws.SendMessage(c.conn, messages.MessageTypeDirectoryStructure, dirStructMsg)

	go func() {
		directoryArr := make([]messages.DirectoryEntry, 1)
		directoryArr[0] = messages.DirectoryEntry{
			Name:               path.Base(startPath),
			RelativePath:       startRelPath,
			ParentRelativePath: path.Dir(startRelPath),
		}
		fileArr := make([]messages.FileEntry, 0)
		collectDirectoryStructure(c.rootDir, startPath, &directoryArr, &fileArr, true)
		dirStruct := messages.DirectoryStructureMessage{
			DirectoryEntries: directoryArr,
			FileEntries:      fileArr,
		}
		dirStructMsg, _ := json.Marshal(dirStruct)
		ws.SendMessage(c.conn, messages.MessageTypeDirectoryStructure, dirStructMsg)
	}()
}

func collectDirectoryStructure(rootPath string, startPath string, directoryArr *[]messages.DirectoryEntry, fileArr *[]messages.FileEntry, extractThumbnail bool) {
	parentRelPath, _ := filepath.Rel(rootPath, startPath)
	dirEntries, err := os.ReadDir(startPath)
	if err != nil {
		log.Fatalf("readDir %s failed, err: %v", startPath, err)
		return
	}
	log.Printf("parsing %s\n", parentRelPath)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			*directoryArr = append(*directoryArr, messages.DirectoryEntry{
				Name:               dirEntry.Name(),
				RelativePath:       filepath.Join(parentRelPath, dirEntry.Name()),
				ParentRelativePath: parentRelPath,
			})
			collectDirectoryStructure(rootPath, filepath.Join(startPath, dirEntry.Name()), directoryArr, fileArr, extractThumbnail)
		} else {
			fileInfo, _ := dirEntry.Info()
			thumbnail := ""
			thumbnailWidth := 1024
			thumbnailHeight := 1024
			if extractThumbnail && isImage(dirEntry.Name()) {
				handler, _ := os.Open(filepath.Join(startPath, dirEntry.Name()))
				srcImage, _, _ := image.Decode(handler)
				originalHeight := srcImage.Bounds().Max.Y
				originalWidth := srcImage.Bounds().Max.X
				thumbnailWidth = 1024 * originalWidth / originalHeight
				dstImage := imaging.Resize(srcImage, thumbnailWidth, 1024, imaging.Lanczos)
				var buf bytes.Buffer
				writer := bufio.NewWriter(&buf)
				jpeg.Encode(writer, dstImage, &jpeg.Options{})
				jpgImage := buf.Bytes()
				thumbnail = fmt.Sprintf("data:image/jpeg;base64,%s", base64.RawStdEncoding.EncodeToString(jpgImage))
				log.Printf("thumbnail generated for %s\n", filepath.Join(startPath, dirEntry.Name()))
			}
			*fileArr = append(*fileArr, messages.FileEntry{
				Name:               dirEntry.Name(),
				RelativePath:       filepath.Join(parentRelPath, dirEntry.Name()),
				ParentRelativePath: parentRelPath,
				FileSize:           uint(fileInfo.Size()),
				Thumbnail:          thumbnail,
				ThumbnailHeight:    thumbnailHeight,
				ThumbnailWidth:     thumbnailWidth,
			})
		}
	}
}

var IMAGE_SUPPORTED = []string{
	".jpeg", ".jpg", ".png", ".gif", ".webp",
}

func isImage(name string) bool {
	return pie.Contains(IMAGE_SUPPORTED, strings.ToLower(filepath.Ext((name))))
}
