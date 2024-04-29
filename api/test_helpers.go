package api

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/db"
	"log"
	"testing"
)

const (
	testmongouri = "mongodb://mongoadmin:bdung@localhost:27017"
)

type testdb struct {
	client *mongo.Client
	db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testmongouri))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
