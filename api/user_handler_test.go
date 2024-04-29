package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "a@asd.com",
		FirstName: "James",
		LastName:  "asd",
		Password:  "sdfdsfdsfffd",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Error("expected user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expected EncryptedPassword to be included in json")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("Expected firstname %s, but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected lastname %s, but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("Expected email %s, but got %s", params.Email, user.Email)
	}
}
