package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DBNAME_HOTEL = "hotel-reservation"
	TestDBNAME   = "hotel-reservation-test"
	DBURI        = "mongodb://mongoadmin:bdung@localhost:27017"
)

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}
