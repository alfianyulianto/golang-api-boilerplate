package main

import (
	"fmt"
	"github.com/alfianyulianto/golang-wilayah-indonesia/wilayah"

	"github.com/alfianyulianto/pds-service/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	username := viperConfig.GetString("database.username")
	password := viperConfig.GetString("database.password")
	host := viperConfig.GetString("database.host")
	port := viperConfig.GetInt("database.port")
	database := viperConfig.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		username, password, host, port, database)
	db := wilayah.ConnectDB("mysql", dsn)
	defer db.Close()

	wilayah.RunMigration(db)
	wilayah.Seed(db, "data")
}
