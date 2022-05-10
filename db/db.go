package db

import (
	"github.com/go-pg/pg/v10"
	"houze_ops_backend/config"
	"log"
	"os"
)

var (
	connectDB *pg.DB
)

func InitConnectionDB() *pg.DB {
	env := config.GetEnvValue()
	connectDB = pg.Connect(&pg.Options{
		User:     env.Database.User,
		Password: env.Database.Password,
		Addr:     env.Database.Host + ":" + env.Database.Port,
		Database: env.Database.Name,
	})
	if connectDB == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	} else {
		log.Printf("Connected to db")
	}
	return connectDB
}

// return database connection
func GetConnectionDB() *pg.DB {
	return connectDB
}
