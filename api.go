package main

import "github.com/gin-gonic/gin"

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

func (api *APIServer) Start() {
	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.Static("/static", "./public")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

	err := r.Run(api.listenAddr)
	if err != nil {
		panic(err)
	}
}
