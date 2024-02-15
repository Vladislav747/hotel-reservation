package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotel-reservation/types"
)

const roomsCollectionName = "rooms"

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, dbname string) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(dbname).Collection(roomsCollectionName),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = resp.InsertedID.(primitive.ObjectID)

	//update the hotel with this room id
	//FIXME Ошибка тут
	//filter := bson.M{"_id": room.HotelID}
	//update := bson.M{"$push": bson.M{"rooms": room.ID}}
	//
	//if err := s.HotelStore.Update(ctx, filter, update); err != nil {
	//	return nil, err
	//}

	return room, nil
}
