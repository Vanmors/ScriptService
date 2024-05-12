package main

import (
	"ScriptService/internal/app"
	"log"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
