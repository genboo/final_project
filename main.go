package main

import (
	"log"
	"os"

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
	gin.SetMode(os.Getenv("gin_mode"))
	server := gin.Default()
	API.Init(server)
	err := server.Run(":80")
	if err != nil {
		log.Println(err)
	}
}
