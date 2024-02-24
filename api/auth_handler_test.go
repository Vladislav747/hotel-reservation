package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "james",
		LastName:  "foo",
		Email:     "james@foo.com",
		Password:  "supersecurepass",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func TestAuthenticate(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb.UserStore)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepass",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var authResp AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Log(err)
	}
}
