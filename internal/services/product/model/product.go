package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price" binding:"required,gt=0"`
	Quantity    int32              `bson:"quantity" json:"quantity" binding:"required,gte=0"`
	Category    string             `bson:"category" json:"category" binding:"required"`
	SKU         string             `bson:"sku" json:"sku" binding:"required"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int32   `json:"quantity" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity    int32   `json:"quantity" binding:"omitempty,gte=0"`
	Category    string  `json:"category"`
}
