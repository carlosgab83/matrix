package main

import (
	"fmt"

	"github.com/carlosgab83/matrix/go/internal/trinity/platform/entrypoint"
)

func main() {
	fmt.Println("Trinity price ingestor starting...")

	myApp, err := entrypoint.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}

	myApp.Run()
	defer myApp.Logger.Close()
}
