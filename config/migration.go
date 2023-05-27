package config

import (
	"api-cariprice/app/entity"
	"fmt"
	"log"
)

func Migration(db Database) {
	fmt.Println("Process Migrating...")

	err := db.DB.AutoMigrate(
		&entity.Cars{},
		&entity.Comparasion{},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database migrated successfully...")
}
