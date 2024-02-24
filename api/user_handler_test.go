package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"net/http/httptest"
	"testing"
)

const (
	dbname       = "hotel-reservation-test"
	testmongouri = "mongodb://mongoadmin:bdung@localhost:27017"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(_ *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testmongouri))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
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
