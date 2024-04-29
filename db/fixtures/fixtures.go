package fixtures

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"
)

func AddBooking(store db.Store, userID, roomID primitive.ObjectID, from, tillDate time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: tillDate,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("UserID %s - AddedBooking\n", booking.UserID)
	return insertedBooking

}

func AddRoom(store db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("SIZE %s -> hotelID %s - AddedRoom\n", room.Size, room.HotelID)
	return insertedRoom
}

func AddHotel(store db.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomsIDS = rooms
	if rooms == nil {
		roomsIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomsIDS,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s -> rooms %s - AddedHotel\n", hotel.Name, hotel.Rooms)
	return insertedHotel
}

func AddUser(store db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s -> %s - AddedUser\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser

}
