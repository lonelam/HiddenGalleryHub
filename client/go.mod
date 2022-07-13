module HiddenGalleryHub/client

go 1.18

require HiddenGalleryHub/common v0.0.0-00010101000000-000000000000

replace HiddenGalleryHub/common => ../common

require github.com/gorilla/websocket v1.5.0

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	golang.org/x/image v0.0.0-20220617043117-41969df76e82 // indirect
)
