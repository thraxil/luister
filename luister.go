package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	const addr = "postgresql://luister@localhost:26257/luister?sslmode=disable"
	db, err := gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&Artist{})
	db.AutoMigrate(&Year{})
	db.AutoMigrate(&Album{})
	db.AutoMigrate(&Song{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Play{})
	db.AutoMigrate(&Rating{})
	db.AutoMigrate(&Tag{})

	fmt.Println("migrated")
}
