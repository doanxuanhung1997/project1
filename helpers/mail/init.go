package mail

import (
	"go.mongodb.org/mongo-driver/mongo"
	"sandexcare_backend/db"
)

/*Resource enum*/
var ()

/*Collection Db*/
func collection(c string) *mongo.Collection {
	return db.Collection(c)
}
