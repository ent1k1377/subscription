package main

import "github.com/ent1k1377/subscriptions/internal/app"

// @title			Subscription API
// @version		1.0
// @description	API for managing user subscriptions
// @host			localhost:8080
// @BasePath		/api
// @schemes		http
func main() {
	app.New().Run()
}
