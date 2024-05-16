package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/db"
	"hotel-reservation/db/fixtures"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		dbName        = os.Getenv("MONGO_DB_NAME")
	)
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(dbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println(user)
	fixtures.AddHotel(store, "asd", "asd", 3, nil)
	hotel := fixtures.AddHotel(store, "Bellucia", "France", 3, nil)
	room := fixtures.AddRoom(store, "large", true, 89.99, hotel.ID)
	booking := fixtures.AddBooking(store, hotel.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println(booking)
	fixtures.AddUser(store, "admin", "admin", true)
	fixtures.AddRoom(store, "small", true, 89.99, hotel.ID)
	fixtures.AddRoom(store, "medium", true, 89.99, hotel.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
