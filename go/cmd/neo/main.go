package main

import (
	"fmt"

	"github.com/carlosgab83/matrix/go/internal/neo/handler"
)

func main() {
	fmt.Println("Neo price collector starting...")

	myApp, err := handler.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}

	myApp.Run()
	defer myApp.Logger.Close()
}
