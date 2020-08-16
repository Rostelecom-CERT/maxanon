package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (m *MongoDB) Open(url string) error {
	var err error
	m.client, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return err
	}
	err = m.client.Connect(context.TODO())
	if err != nil {
		return err
	}
	err = m.client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	m.collection = m.client.Database("reputation").Collection("records")
	return nil
}

func (m *MongoDB) Insert(data Data) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.collection.InsertOne(ctx, data)

	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) InsertBulk(bulk []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.collection.InsertMany(ctx, bulk)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Get(ip string) (*Data, error) {
	var result Data
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	type req struct {
		IP string `bson:"ip"`
	}
	err := m.collection.FindOne(ctx, req{IP: ip}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return &result, nil
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MongoDB) Exist(collName string) (bool, error) {
	result, err := m.client.Database("reputation").ListCollectionNames(context.TODO(), bson.D{{"options.capped", true}})
	if err != nil {
		return false, err
	}

	for _, coll := range result {
		if coll == collName {
			return true, nil
		}
	}
	return false, nil
}
