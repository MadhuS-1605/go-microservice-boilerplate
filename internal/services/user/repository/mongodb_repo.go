package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/services/user/model"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *database.MongoDB) UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *mongoUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepository) Update(ctx context.Context, id string, user *model.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *mongoUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *mongoUserRepository) List(ctx context.Context, page, limit int, search string) ([]*model.User, int64, error) {
	filter := bson.M{}

	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip
	skip := (page - 1) * limit

	// Find documents with pagination
	cursor, err := r.collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &[]int64{int64(skip)}[0],
		Limit: &[]int64{int64(limit)}[0],
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
