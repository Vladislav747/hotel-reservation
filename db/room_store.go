package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"hotel-reservation/types"
)

const roomsCollectionName = "rooms"

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(roomsCollectionName),
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	_, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	//room.ID = resp.InsertedID.(primitive.ObjectID)

	////update the hotel with this room id
	////FIXME Ошибка тут
	//filter := Map{"_id": room.HotelID}
	//update := Map{"$push": bson.M{"rooms": room.ID}}
	////
	//if err := s.HotelStore.Update(ctx, filter, update); err != nil {
	//	return nil, err
	//}

	return room, nil
}
