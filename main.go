package main

import (
	"gcp/config"
	"gcp/services"
)

func main() {
	// Connect to the database
	config.ConnectDatabase()

	// Trigger region sync
	// if err := services.FetchAndStoreRegions(); err != nil {
	// 	fmt.Println("Error syncing regions:", err)
	// }

	services.FetchAndInsertSkus()
}
