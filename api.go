package main

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
