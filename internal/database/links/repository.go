package links

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

const collection = "links"

func New(db *mongo.Database, timeout time.Duration) *Repository {
	return &Repository{db: db, timeout: timeout}
}

type Repository struct {
	db      *mongo.Database
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateReq) (database.Link, error) {
	var l database.Link
	// implemented
	now := time.Now()
	l = database.Link{
		ID:        req.ID,
		Title:     req.Title,
		URL:       req.URL,
		Images:    req.Images,
		Tags:      req.Tags,
		UserID:    req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := r.db.Collection(collection).InsertOne(ctx, l)
	if err != nil {
		return database.Link{}, fmt.Errorf("mongo InsertOne: %w", err)
	}

	return l, nil
}

func (r *Repository) FindByUserAndURL(ctx context.Context, link, userID string) (database.Link, error) {
	var l database.Link
	// implemented
	result := r.db.Collection(collection).FindOne(ctx, bson.M{"url": link, "userID": userID})
	if err := result.Decode(&l); err != nil {
		return l, fmt.Errorf("mongo Decode: %w", err)
	}
	return l, nil
}

func (r *Repository) FindByCriteria(ctx context.Context, criteria Criteria) ([]database.Link, error) {
	//implement me
	var links []database.Link

	findFilter := bson.M{}
	if criteria.UserID != nil {
		findFilter["userID"] = *criteria.UserID
	}
	if len(criteria.Tags) > 0 {
		//TODO
		findFilter["tags"] = bson.M{"$in": criteria.Tags}
	}

	opts := options.Find()
	if criteria.Limit != nil {
		opts.SetLimit(*criteria.Limit)
	}
	if criteria.Offset != nil {
		opts.SetSkip(*criteria.Offset)
	}

	cursor, err := r.db.Collection(collection).Find(ctx, findFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo Find: %w", err)
	}

	for cursor.Next(ctx) {
		var l database.Link
		if err := cursor.Decode(&l); err != nil {
			return nil, fmt.Errorf("mongo Decode: %w", err)
		}
		links = append(links, l)
	}

	return links, nil

}
