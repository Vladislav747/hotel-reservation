package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
)

func main() {

	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME_HOTEL).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME_HOTEL)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME_HOTEL)

	hotel := types.Hotel{
		Name:     "B",
		Location: "C",
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Type:  types.SingleRoomType,
			Price: 99,
		},
		{
			Type:  types.DeluxeRoomType,
			Price: 199,
		},
		{
			Type:  types.SeaSideRoomType,
			Price: 299,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

}
