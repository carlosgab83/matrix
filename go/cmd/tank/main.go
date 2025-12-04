package main

import (
	"fmt"

	"github.com/carlosgab83/matrix/go/internal/tank/handler"
)

func main() {
	fmt.Println("Tank notifier starting...")

	myApp, err := handler.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}

	myApp.Run()
	defer myApp.Logger.Close()
}
