package main

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"HiddenGalleryHub/server/routes"
	"flag"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var flagPort = flag.String("port", "5555", "the hosting port")

func main() {
	db, err := gorm.Open(sqlite.Open("../gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	models.Migrate(db)
	models.Init(db)
	pool := connections.CreatePool()
	routes.RunWithWebsocketUpgrader("localhost:"+*flagPort, pool, db)

	<-pool.Done
}
