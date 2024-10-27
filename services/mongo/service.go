// Copyright © ivanlobanov. All rights reserved.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *mongo.Client {
	connLine := fmt.Sprintf(`mongodb://%s:%d/main?connectTimeoutMS=5000&socketTimeoutMS=10000&maxPoolSize=30&minPoolSize=0&maxConnecting=3`,
		viper.GetString("mongo.host"), viper.GetInt("mongo.port"))
	client, err := mongo.Connect(options.Client().ApplyURI(connLine))
	if err != nil {
		logger.Fatal(fmt.Sprintf("error MongoDB connect: %v", err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 3; i++ {
		err = client.Ping(ctx, nil)
		if err == nil {
			break
		}
		if i == 2 {
			logger.Fatal("MongoDB is not respond")
		}
		logger.Warn(fmt.Sprintf("error while ping to postgresql: %v", err))
		time.Sleep(3 * time.Second)
	}

	logger.Info("succesful connection to MongoDB")
	return client
}
