package users

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sethvargo/go-envconfig"
	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/env/config"
	"log"
	"os"
	"testing"
	"time"
)

var (
	usersRepo *Repository
	conn      *pgx.Conn
	userID    uuid.UUID
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	userID = uuid.New()
	var cfg config.Config
	ctx := context.Background()

	if err := envconfig.Process(ctx, &cfg); err != nil { //nolint:typecheck
		log.Fatalf("env processing: %v", err)
	}

	usersConn, err := pgx.Connect(ctx, cfg.UsersDB.ConnectionURL())
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	conn = usersConn

	usersRepo = New(usersConn, 5*time.Second)
}

func shutdown() {
	ctx := context.Background()
	_, _ = usersRepo.userDB.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	_ = conn.Close(ctx)
}

func TestRepository_Create_Find(t *testing.T) {
	ctx := context.Background()

	username := "test_username"

	created, err := usersRepo.Create(
		ctx, CreateUserReq{
			ID:       userID,
			Username: username,
			Password: "test_password",
		},
	)

	handleError(t, err)
	if created.ID != userID {
		t.Fatal("create id is incorrect: ", created.ID)
	}

	findByID, err := usersRepo.FindByID(ctx, userID)
	handleError(t, err)

	if findByID.Username != username || findByID.ID != userID {
		t.Fatal("found by username is incorrect: ", created.Username, created.ID)
	}

	findByUsername, err := usersRepo.FindByUsername(ctx, username)
	handleError(t, err)

	if findByUsername.Username != username || findByUsername.ID != userID {
		t.Fatal("found by username is incorrect: ", created.Username, created.ID)
	}
}

func handleError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
