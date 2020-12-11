package main

import (
	"context"
	gz_mongo "github.com/gzshanyu/go-utils/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func main() {
	var (
		err        error
		mongoCli   *mongo.Client
		db         *mongo.Database
		collection *mongo.Collection
		result     *mongo.SingleResult
	)

	if mongoCli, err = gz_mongo.NewConnect([]string{"139.155.56.188:27017"}, 5*time.Second); err != nil {
		log.Println(err)
		return
	}
	// 选择数据库
	db = mongoCli.Database("gzshanyu")
	// 选择表
	collection = db.Collection("cron_logs")

	result = collection.FindOne(context.TODO(), map[string]interface{}{"name": "jobs"})
	var infter = make(map[string]interface{})
	result.Decode(&infter)
	log.Println(infter)
}
