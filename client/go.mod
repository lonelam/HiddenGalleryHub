module HiddenGalleryHub/client

go 1.18

require HiddenGalleryHub/common v0.0.0-00010101000000-000000000000

replace HiddenGalleryHub/common => ../common

require github.com/gorilla/websocket v1.5.0
