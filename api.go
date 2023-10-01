package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	listenAddr string
	store      Store
}

func NewAPIServer() *APIServer {
	env := GetEnv()

	store, err := NewStoreInstance()
	if err != nil {
		panic(err)
	}

	return &APIServer{
		listenAddr: env.Addr,
		store:      store,
	}
}

func (api *APIServer) handleIndex(c *gin.Context) {
	docs, err := api.store.GetPages()
	if err != nil {
		docs = []PageDoc{}
		log.Default().Println(err)
	}

	// Reverse the docs so that the newest is first.
	for i, j := 0, len(docs)-1; i < j; i, j = i+1, j-1 {
		docs[i], docs[j] = docs[j], docs[i]
	}

	c.HTML(200, "index.html", gin.H{
		"docs": docs[:5],
	})
}

func (api *APIServer) handleCreate(c *gin.Context) {
	c.HTML(200, "create.html", gin.H{})
}

func (api *APIServer) handleCreatePage(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	content := c.PostForm("content")
	date := time.Now().UTC().Format(time.RFC3339)

	hash := ConvertTitleToHash(title)

	err := api.store.CreatePage(PageDocRequest{
		Hash:    hash,
		Name:    title,
		Content: content,
		Date:    date,
		Author:  author,
	})

	if err != nil {
		log.Default().Println("Error creating page:", err)
		c.String(500, "Error creating page")
		return
	}

	c.String(200, fmt.Sprintf("Created <a href=\"/%s\">%s</a> successfully!", hash, title))
}

func (api *APIServer) Start() {
	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.Static("/static", "./public")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", api.handleIndex)
	r.GET("/new", api.handleCreate)

	r.POST("/create", api.handleCreatePage)

	err := r.Run(api.listenAddr)
	if err != nil {
		panic(err)
	}
}
