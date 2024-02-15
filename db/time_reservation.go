package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"hotel-reservation/types"
)

const timeReservationCollectionName = "rooms"

type TimeReservationStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoTimeReservationStore struct {
	coll *mongo.Collection

	HotelStore
}

func NewMongoTimeReservationStore(client *mongo.Client, dbname string) *MongoTimeReservationStore {
	return &MongoTimeReservationStore{
		coll: client.Database(dbname).Collection(timeReservationCollectionName),
	}
}

func (s *MongoRoomStore) InsertReservation(ctx context.Context, room *types.Room) (*types.Room, error) {
	//Сначала отсортировать - использовать индекс - блочить на уровне транзакции
	//resp, err := s.coll.FindOne(ctx, room)
	//var time_reservation types.TimeReservation
	//if err := s.coll.FindOne(ctx, bson.M{"_id": ToObjectID(room.ID)}).Decode(&user); err != nil {
	//	return nil, err
	//}
	//
	//
	//if err!= nil {
	//    return nil, err
	//}
	//
	//resp, err = s.coll.InsertOne(ctx, room)
	//if err != nil {
	//	return nil, err
	//}
	//
	//room.ID = resp.InsertedID.(primitive.ObjectID)
	//
	////update the hotel with this room id
	//filter := bson.M{"_id": room.HotelID}
	//update := bson.M{"$push": bson.M{"rooms": room.ID}}
	//
	//if err := s.HotelStore.Update(ctx, filter, update); err != nil {
	//	return nil, err
	//}

	return room, nil
}
