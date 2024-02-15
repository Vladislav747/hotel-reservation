package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeReservation struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID    int                 `bson:"room_id" json:"room_id"`
	StartDate primitive.Timestamp `bson:"start_date" json:"start_date"`
	EndDate   primitive.Timestamp `bson:"end_date" json:"end_date"`
	UserID    int                 `bson:"user_id" json:"user_id"`
}
