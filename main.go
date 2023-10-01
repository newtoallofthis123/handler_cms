package main

func main() {
	api := NewAPIServer()
	api.store.Init()
	api.Start()
}
