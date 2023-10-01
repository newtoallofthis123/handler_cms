package main

import (
	"fmt"
)

func main() {
	api := NewAPIServer()
	api.store.Init()

	fmt.Println("Hello, world!")
}
