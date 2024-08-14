package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "mongodb://127.0.0.1:27017/?replicaSet=rs0"
// "mongodb://127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/?replicaSet=rs0"

type (
	User struct {
		Name string `bson:"name" json:"name"`
		Age  int    `bson:"age" json:"age"`
	}

	UserWithID struct {
		ID   uint64 `bson:"_id" json:"id"` // primitive.ObjectID
		Name string `bson:"name" json:"name"`
		Age  int    `bson:"age" json:"age"`
	}
)

func main() {
	dbHost := os.Getenv("DB_HOST")

	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(dbHost)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	// collection := client.Database("testdb").Collection("users")

	// // INSERT
	// user := User{"Alex", 50}

	// insertResult, err := collection.InsertOne(ctx, user)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Added document with ID:", insertResult.InsertedID)

	// // FIND
	// var result User
	// filter := bson.M{"name": "Alex"}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Found document:", result)

	// // UPDATE
	// update := bson.M{"$set": bson.M{"age": 45}}
	// updateResult, err := collection.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Updated document: %v\n", updateResult.ModifiedCount)

	// // DELETE
	// deleteResult, err := collection.DeleteOne(ctx, filter)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Deleted document: %v\n", deleteResult.DeletedCount)

	// collectionMany := client.Database("testdb").Collection("usersWithCustomID")

	// // INSERT MANY
	// users := []interface{}{
	// 	UserWithID{1, "Jack", 25},
	// 	UserWithID{2, "Reaper", 35},
	// }

	// insertManyResult, err := collectionMany.InsertMany(ctx, users)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Inserted IDs:", insertManyResult.InsertedIDs)

	// // FIND MANY
	// cur, err := collectionMany.Find(ctx, bson.D{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cur.Close(ctx)

	// for cur.Next(ctx) {
	// 	var elem UserWithID
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Document:", elem)
	// }

	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

}
