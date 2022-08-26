package main

import (
	"FunNow/url-shortener/constants"
	"FunNow/url-shortener/routers"
	"math/rand"
	"time"
)

func main() {
	// Initialize rand seed when the server is up
	rand.Seed(time.Now().UnixNano())

	// Run server
	engine := routers.InitRouter()
	panic(engine.Run(":" + constants.WebServerPort))
}
