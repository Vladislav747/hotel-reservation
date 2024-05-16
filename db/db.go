package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pagination struct {
	Limit int64
	Page  int64
}

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
