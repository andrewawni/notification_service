package main

import (
	"log"
)

func main() {
	app := App{}
	app.Init()
	log.Printf("server running")
	app.Run(":8000")
}
