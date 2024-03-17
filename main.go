package main

import (
	"log"
	"os"
	"strconv"

	"github.com/genboo/final_project/api"
	"github.com/genboo/final_project/cache"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	DefaultCacheDir = "image_cache"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	capacity := os.Getenv("max_cache_capacity")
	if capacity == "" {
		log.Fatalln("необходимо задать параметр max_cache_capacity")
	}
	val, err := strconv.Atoi(capacity)
	if err != nil {
		log.Fatalln(err)
	}
	gin.SetMode("release")
	server := gin.Default()
	cache.Init(DefaultCacheDir, val)
	api.Init(server)
	err = server.Run(":80")
	if err != nil {
		log.Println(err)
	}
}
