package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This creates a new Store Instance that implements the Store interface.
func NewStoreInstance() (*DBInstance, error) {
	env := GetEnv()

	// We get all the info we need from the environment variables.
	// The context.TODO() is a placeholder for a real context.
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.URI))
	if err != nil {
		return nil, err
	}

	// This is a ping to the database to make sure we can connect
	// If we can't connect, we return an error.
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

// We initialize the cache by hydrating it with the data from the database.
func (db *DBInstance) Init() {
	db.HydrateCache()
}

// This is a helper function to get all the pages from the database.
// Internally, it uses the mongodb client to get all the pages.
// but stores them in a slice of PageDoc structs in the cache.
// This is so we don't have to query the database every time we want to get a page.
func (db *DBInstance) getPagesFromDB() ([]PageDoc, error) {
	cursor, err := db.mongodb.Collection("page").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var pages []PageDoc

	// Recursively decode the cursor into a slice of PageDoc structs.
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

// This is a helper function to hydrate the cache.
// This void function gets called every time we create, update, or delete a page.
// All read operations are cached to save on database queries.
func (db *DBInstance) HydrateCache() {
	pages, err := db.getPagesFromDB()
	if err != nil {
		log.Default().Println("Error hydrating cache:", err)
		return
	}

	db.cache.docs = pages
}

// Add a page to the database and hydrates the cache.
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

// Get a page from the cache.
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

// Get all the pages from the cache.
func (db *DBInstance) GetPages() ([]PageDoc, error) {
	return db.cache.docs, nil
}

// Update a page in the database and hydrates the cache.
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

// Delete a page from the database and hydrates the cache.
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

// Search the cache
func (db *DBInstance) SearchPages(query string) ([]PageDoc, error) {
	query = strings.ToLower(query)
	docs, err := db.GetPages()
	if err != nil {
		log.Default().Println("Error hydrating cache:", err)
		return nil, err
	}

	var results []PageDoc

	for _, doc := range docs {
		docString := fmt.Sprintf("%v", doc)
		docString = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(docString, "-", ""), ":", ""), "{", ""))
		if strings.Contains(docString, query) {
			results = append(results, doc)
		}
	}

	return results, nil
}
