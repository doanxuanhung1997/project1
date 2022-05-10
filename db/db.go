package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"houze_ops_backend/config"
	"sync"
)

var (
	db     *mongo.Database
	client *mongo.Client
	ct     = context.Background()
	err    error
	once   sync.Once
)

func GetContext() context.Context {
	return ct
}

func GetClient() *mongo.Client {
	return client
}

func GetDatabase() *mongo.Database {
	return db
}

func InitDb() error {
	once.Do(func() {
		env := config.GetEnvValue()
		ctx := context.Background()
		// Options to the database.
		clientOpts := options.Client().ApplyURI(env.Database.Host)
		client, err := mongo.Connect(ctx, clientOpts)
		if err != nil {
			panic(err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		db = client.Database(env.Database.Name)
		fmt.Println("database name:  " + db.Name())
	})
	return err
}

// Collection returns database
func Collection(collection string) *mongo.Collection {
	return db.Collection(collection)
}
