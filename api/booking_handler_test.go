package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db/fixtures"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserBookingHandler(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user           = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(tdb.Store, "bar hotel", "a", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("non 200 status Code got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookingResp, "bookingResp")
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user      = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(tdb.Store, "bar hotel", "a", 4, nil)
		room      = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)

		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	fmt.Println(booking, "booking")
	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookings, "bookings")
	if len(bookings) != 1 {
		t.Fatalf("expecting 1 booking, got %d", len(bookings))
	}
	have := bookings[0]
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s, got %s", booking.UserID, have.UserID)
	}
	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("exprected a non 200 status Code got %d", resp.StatusCode)
	}

}
