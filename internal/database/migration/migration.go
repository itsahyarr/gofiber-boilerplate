package migration

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"

	"github.com/itsahyarr/gofiber-boilerplate/pkg/database"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/logger"
)

// RunMigrations executes all database migrations (indexes, etc.)
func RunMigrations(db *database.MongoDB) {
	logger.Info("Running database migrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. User Indexes
	migrateUserIndexes(ctx, db)

	// Add more migration modules here as needed

	logger.Info("Database migrations completed successfully")
}

func migrateUserIndexes(ctx context.Context, db *database.MongoDB) {
	collection := db.Collection("users")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		logger.Error("Failed to create user indexes", zap.Error(err))
	} else {
		logger.Info("User indexes verified/created")
	}
}
