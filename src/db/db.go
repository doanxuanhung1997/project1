package db

import (
	"github.com/go-pg/pg/v10"
	"houze_ops_backend/configs"
	"log"
	"os"
)

var (
	connectDB *pg.DB
)

func InitConnectionDB() *pg.DB {
	env := configs.GetEnvConfig()
	connectDB = pg.Connect(&pg.Options{
		User:     env.DBUser,
		Password: env.DBPassword,
		Addr:     env.DBHost + ":" + env.DBPort,
		Database: env.DBName,
	})
	if connectDB == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	} else {
		log.Printf("Connected to db")
	}
	return connectDB
}

// GetConnectionDB return database connection
func GetConnectionDB() *pg.DB {
	return connectDB
}
