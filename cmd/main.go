package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	api "hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db"
	"hotel-reservation/errors"
	"log"
	"os"
	"os/signal"
)

var config = fiber.Config{
	ErrorHandler: errors.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore = db.NewMongoUserStore(client)

		roomStore    = db.NewMongoRoomStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api/")
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		// api/v1/admin
		admin = apiv1.Group("/admin", middleware.AdminAuth)
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
	auth.Post("/auth", authHandler.HandleAuthenticate)

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

	//room handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Get("/room/:id/book", roomHandler.HandleBookRoom)

	//booking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	//admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	err = app.Listen(listenAddr)
	if err != nil {
		fmt.Printf("Listen Connection error: %s\n", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
