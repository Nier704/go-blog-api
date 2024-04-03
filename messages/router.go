package messages

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Author string             `bson:"author"`
	Text   string             `bson:"text"`
	Date   string             `bson:"date"`
}

type MessageDTO struct {
	Author string `bson:"author"`
	Text   string `bson:"text"`
	Date   string `bson:"date"`
}

type MessageList []Message

var db *mongo.Database

func Router(r chi.Router) {
	if err := openDBConn(); err != nil {
		panic(err)
	}

	r.Get("/", allMessages)
	r.Post("/", sendMessage)
}

func allMessages(w http.ResponseWriter, r *http.Request) {
	coll := db.Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var msg bson.M
		if err = cursor.Decode(&msg); err != nil {
			panic(err)
		}
		fmt.Println(msg)
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {

}

func openDBConn() (err error) {
	godotenv.Load()
	db_password := os.Getenv("DB_PASSWORD")

	db_url := fmt.Sprintf("mongodb+srv://lorecode12:%s@cluster0.qw8e8tb.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", db_password)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db_url))
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	db = client.Database("blog-go")
	fmt.Println("database connection established")
	return nil
}
