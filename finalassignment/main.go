package main

import (
	"finalassignment.id/finalassignment/database"
	_ "finalassignment.id/finalassignment/docs"
	"finalassignment.id/finalassignment/routers"
)

// @title           Final Assignment
// @version         1.0
// @description     Final Assignment

// @contact.name   zulkarnaen
// @contact.email  premiumforspot@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	database.StartDB()
	var port = ":8080"
	routers.StartServer().Run(port)
}
