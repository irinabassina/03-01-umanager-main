package links

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/env/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

var (
	linksRepo *Repository
	client    *mongo.Client
	userID    string
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	userID = "test_uuid"
	var cfg config.Config
	ctx := context.Background()

	if err := envconfig.Process(ctx, &cfg); err != nil { //nolint:typecheck
		log.Fatalf("env processing: %v", err)
	}

	linksDB, err := mongo.Connect(
		ctx, &options.ClientOptions{
			ConnectTimeout: &cfg.LinksDB.ConnectTimeout,
			Hosts:          []string{fmt.Sprintf("%s:%d", cfg.LinksDB.Host, cfg.LinksDB.Port)},
			Auth: &options.Credential{
				AuthMechanism: "SCRAM-SHA-256",
				AuthSource:    cfg.LinksDB.Name,
				Username:      cfg.LinksDB.User,
				Password:      cfg.LinksDB.Password,
			},
			MaxPoolSize: &cfg.LinksDB.MaxPoolSize,
			MinPoolSize: &cfg.LinksDB.MinPoolSize,
		},
	)
	if err != nil {
		log.Fatalf("mongo.Connect: %v", err)
	}
	client = linksDB
	linksRepo = New(linksDB.Database(cfg.LinksDB.Name), 5*time.Second)

}

func shutdown() {
	ctx := context.Background()

	_, _ = linksRepo.db.Collection(collection).DeleteMany(ctx, bson.M{"userID": userID})

	_ = client.Disconnect(context.Background())
}

func TestRepository_Create_Find_By_User_URL(t *testing.T) {
	ctx := context.Background()

	id := primitive.NewObjectID()
	url := "https://ya.ru"
	created, err := linksRepo.Create(
		ctx, CreateReq{
			ID:     id,
			URL:    url,
			Title:  "ya main page",
			Tags:   []string{"search", "yandex"},
			Images: []string{},
			UserID: userID, // created user id
		},
	)
	handleError(t, err)
	if created.ID != id {
		t.Fatal("create id is incorrect: ", created.ID)
	}

	found, err := linksRepo.FindByUserAndURL(ctx, url, userID)
	handleError(t, err)
	if found.UserID != userID || found.URL != url {
		t.Fatal("found id is incorrect: ", found.ID, found.URL)
	}
}

func TestRepository_Find_By_Criteria(t *testing.T) {
	ctx := context.Background()

	_, err := linksRepo.Create(
		ctx, CreateReq{
			ID:     primitive.NewObjectID(),
			URL:    "https://ya.ru",
			Title:  "ya main page",
			Tags:   []string{"search", "yandex_find"},
			Images: []string{},
			UserID: userID, // created user id
		},
	)
	handleError(t, err)

	_, err = linksRepo.Create(
		ctx, CreateReq{
			ID:     primitive.NewObjectID(),
			URL:    "https://google.com",
			Title:  "google main page",
			Tags:   []string{"search", "google_find"},
			Images: []string{},
			UserID: userID, // created user id
		},
	)
	handleError(t, err)

	found, err := linksRepo.FindByCriteria(ctx, Criteria{
		Tags: []string{"yandex_find", "google_find"},
	})
	handleError(t, err)
	if len(found) != 2 {
		t.Fatal("found is incorrect: ", found)
	}

}

func handleError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
