package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/bin/api"
	"hotel-reservation/bin/api/middleware"
	"hotel-reservation/db"
	"log"
	"os"
	"os/signal"
)

const dburi = "mongodb://mongoadmin:bdung@localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":3001", "The listen address of the API server")
	flag.Parse()

	var (
		userStore   = db.NewMongoUserStore(client)
		userHandler = api.NewUserHandler(userStore)
		roomStore   = db.NewMongoRoomStore(client)
		hotelStore  = db.NewMongoHotelStore(client)
		store       = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)

		app   = fiber.New(config)
		auth  = app.Group("/api/")
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	//Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	fmt.Println("App Starting")

	//auth
	auth.Post("/auth/:id", authHandler.HandleAuthenticate)

	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*listenAddr)
}
