package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
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

func (api *APIServer) handleSearch(c *gin.Context) {
	query := c.Query("q")

	docs, err := api.store.SearchPages(query)
	if err != nil {
		docs = []PageDoc{}
		log.Default().Println(err)
	}

	c.HTML(200, "search.html", gin.H{
		"num":     len(docs),
		"results": docs,
	})
}

func (api *APIServer) handleGetPage(c *gin.Context) {
	hash := c.Param("hash")

	page, err := api.store.GetPage(hash)
	if err != nil {
		log.Default().Println("Error getting page:", err)
		c.String(500, "Error getting page")
		return
	}

	html := blackfriday.MarkdownCommon([]byte(page.Content))
	page.Content = string(html)

	c.HTML(200, "page.html", gin.H{
		"Page": page,
	})
}

func (api *APIServer) handleListPage(c *gin.Context) {
	pages, err := api.store.GetPages()
	if err != nil {
		log.Default().Println("Error getting page:", err)
		c.String(500, "Error getting page")
		return
	}

	//reverse the pages so that the newest is first
	for i, j := 0, len(pages)-1; i < j; i, j = i+1, j-1 {
		pages[i], pages[j] = pages[j], pages[i]
	}

	c.HTML(200, "list.html", gin.H{
		"Pages": pages,
	})
}

func (api *APIServer) Start() {
	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.Static("/static", "./public")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", api.handleIndex)
	r.GET("/new", api.handleCreate)
	r.GET("/search", api.handleSearch)
	r.GET("/quips/:hash", api.handleGetPage)
	r.GET("/list", api.handleListPage)

	r.POST("/create", api.handleCreatePage)

	err := r.Run(api.listenAddr)
	if err != nil {
		panic(err)
	}
}
