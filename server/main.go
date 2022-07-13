package main

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"HiddenGalleryHub/server/routes"
	"flag"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var flagPort = flag.Int("port", 5555, "the hosting port")

func main() {
	flag.Parse()
	db, err := gorm.Open(sqlite.Open("../gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	models.Migrate(db)
	models.Init(db)
	pool := connections.CreatePool()
	routes.RunWithWebsocketUpgrader(fmt.Sprintf("localhost:%d", *flagPort), pool, db)

	<-pool.Done
}
