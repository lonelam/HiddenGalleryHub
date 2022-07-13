package main

import (
	"HiddenGalleryHub/server/connections"
	"HiddenGalleryHub/server/models"
	"HiddenGalleryHub/server/routes"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("../gallery.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	models.Migrate(db)
	models.Init(db)
	pool := connections.CreatePool()
	routes.RunWithWebsocketUpgrader("localhost:5555", pool, db)

	// autotls.Run(router, "localhost:5555")

	<-pool.Done
}
