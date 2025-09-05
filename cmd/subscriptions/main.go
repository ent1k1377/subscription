package main

import "github.com/ent1k1377/subscriptions/internal/app"

// @title			Subscription API
// @version		1.0
// @description	API для управления подписками пользователей
// @host			localhost:8080
// @BasePath		/api
// @schemes		http
func main() {
	app.New().Run()
}
