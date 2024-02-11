package main

import (
	"context"
	"fmt"
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

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME_HOTEL)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME_HOTEL)

	hotel := types.Hotel{
		Name:     "B",
		Location: "C",
	}

	room := types.Room{
		Type:  types.SingleRoomType,
		Price: 1,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, &room)

	if err != nil {
		log.Fatal(err)
	}

	room.HotelID = insertedHotel.ID

	fmt.Println(insertedHotel)
	fmt.Println(insertedRoom)
}
