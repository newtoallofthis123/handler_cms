package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Env struct {
	URI  string
	DB   string
	Addr string
}

type DBInstance struct {
	client  *mongo.Client
	mongodb *mongo.Database
	cache   CacheStore
}

type Store interface {
	HydrateCache()
	Init()
	GetPage(hash string) (PageDoc, error)
	CreatePage(PageDocRequest) error
	UpdatePage(PageDocRequest) error
	DeletePage(string) error
	GetPages() ([]PageDoc, error)
}

type PageDoc struct {
	DocID   primitive.ObjectID `bson:"_id,omitempty"`
	Hash    string             `bson:"hash"`
	Name    string             `bson:"name"`
	Content string             `bson:"content"`
	Date    time.Time          `bson:"date"`
	Author  string             `bson:"author"`
}

type CacheStore struct {
	docs []PageDoc
}

type PageDocRequest struct {
	Hash    string `json:"hash"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Date    string `json:"date"`
	Author  string `json:"author"`
}
