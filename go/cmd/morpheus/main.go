package main

import (
	"fmt"

	"github.com/carlosgab83/matrix/go/internal/morpheus/handler"
)

func main() {
	fmt.Println("Morpheus price ingestor starting...")

	myApp, err := handler.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}

	myApp.Run()
	defer myApp.Logger.Close()
}
