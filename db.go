package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStoreInstance() (*DBInstance, error) {
	env := GetEnv()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.URI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &DBInstance{
		client:  client,
		mongodb: client.Database(env.DB),
		cache:   CacheStore{},
	}, nil
}

func (db *DBInstance) Init() {
	db.HydrateCache()
}

func (db *DBInstance) getPages() ([]PageDoc, error) {
	cursor, err := db.mongodb.Collection("page").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var pages []PageDoc

	for cursor.Next(context.Background()) {
		var page PageDoc
		err := cursor.Decode(&page)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (db *DBInstance) HydrateCache() {
	pages, err := db.getPages()
	if err != nil {
		log.Default().Println("Error hydrating cache:", err)
		return
	}

	db.cache.docs = pages
}

func (db *DBInstance) CreatePage(req PageDocRequest) error {
	_, err := db.mongodb.Collection("page").InsertOne(context.Background(), bson.M{
		"hash":    req.Hash,
		"name":    req.Name,
		"content": req.Content,
		"date":    req.Date,
		"author":  req.Author,
	})

	if err != nil {
		return err
	}

	db.HydrateCache()
	return nil
}

func (db *DBInstance) GetPage(hash string) (PageDoc, error) {
	var page PageDoc

	for _, p := range db.cache.docs {
		if p.Hash == hash {
			page = p
			return page, nil
		}
	}

	return PageDoc{}, nil
}

func (db *DBInstance) GetPages() ([]PageDoc, error) {
	return db.cache.docs, nil
}

func (db *DBInstance) UpdatePage(req PageDocRequest) error {
	_, err := db.mongodb.Collection("page").UpdateOne(context.Background(), bson.M{
		"hash": req.Hash,
	}, bson.M{
		"$set": bson.M{
			"name":    req.Name,
			"content": req.Content,
			"date":    req.Date,
			"author":  req.Author,
		},
	})

	if err != nil {
		return err
	}

	db.HydrateCache()
	return nil
}

func (db *DBInstance) DeletePage(hash string) error {
	_, err := db.mongodb.Collection("page").DeleteOne(context.Background(), bson.M{
		"hash": hash,
	})

	if err != nil {
		return err
	}

	db.HydrateCache()
	return nil
}
