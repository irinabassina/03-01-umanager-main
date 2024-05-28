package users

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

func New(userDB *pgx.Conn, timeout time.Duration) *Repository {
	return &Repository{userDB: userDB, timeout: timeout}
}

type Repository struct {
	userDB  *pgx.Conn
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateUserReq) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// implemented
	now := time.Now()
	u = database.User{
		ID:        req.ID,
		Username:  req.Username,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.userDB.Exec(ctx, `INSERT INTO users (id, username, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		u.ID, u.Username, u.Password, now, now)
	if err != nil {
		return u, fmt.Errorf("postgres Exec: %w", err)
	}

	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	// implemented
	err := r.userDB.QueryRow(ctx, `SELECT * FROM users WHERE id=$1`, userID).Scan(
		&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return u, fmt.Errorf("postgres QueryRow Decode: %w", err)
	}

	return u, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	// implemented
	err := r.userDB.QueryRow(ctx, `SELECT * FROM users WHERE username=$1`, username).Scan(
		&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return u, fmt.Errorf("postgres QueryRow Decode: %w", err)
	}
	return u, nil
}
