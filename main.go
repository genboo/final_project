package main

import (
	"log"

	"github.com/genboo/final_project/api"
	"github.com/genboo/final_project/cache"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	gin.SetMode("release")
	server := gin.Default()
	cache.Init("image_cache")
	api.Init(server)
	err := server.Run(":80")
	if err != nil {
		log.Println(err)
	}
}
