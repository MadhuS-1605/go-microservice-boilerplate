package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-microservice-boilerplate/internal/database"
	"go-microservice-boilerplate/internal/services/product/model"
)

type mongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *database.MongoDB) ProductRepository {
	return &mongoProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *mongoProductRepository) Create(ctx context.Context, product *model.Product) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *mongoProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *mongoProductRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var product model.Product
	err := r.collection.FindOne(ctx, bson.M{"sku": sku}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id string, product *model.Product) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	product.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"quantity":    product.Quantity,
			"category":    product.Category,
			"updated_at":  product.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *mongoProductRepository) List(ctx context.Context, page, limit int, search, category string) ([]*model.Product, int64, error) {
	filter := bson.M{}

	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
			{"sku": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if category != "" {
		filter["category"] = bson.M{"$regex": category, "$options": "i"}
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

	var products []*model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Helper functions for BSON type conversion
func getInt64FromBSON(doc bson.M, key string) int64 {
	if val, ok := doc[key]; ok {
		switch v := val.(type) {
		case int64:
			return v
		case int32:
			return int64(v)
		case int:
			return int64(v)
		}
	}
	return 0
}

func getFloat64FromBSON(doc bson.M, key string) float64 {
	if val, ok := doc[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int64:
			return float64(v)
		case int32:
			return float64(v)
		case int:
			return float64(v)
		}
	}
	return 0.0
}
