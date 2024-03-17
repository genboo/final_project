package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/genboo/final_project/cache"
	"github.com/gin-gonic/gin"
)

const (
	paramUrl      = "any"
	paramWidth    = "width"
	paramHeight   = "height"
	contextParams = "params"
)

func Init(server *gin.Engine) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s\n", err)
		}
	}()
	server.GET("/fill/:width/:height/*any", validate, preview)
}

func preview(c *gin.Context) {
	params := c.MustGet(contextParams).(cache.Params)
	image, err := cache.ImageCache.GetImage(c.Request.Header, params)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	_, err = c.Writer.Write(image)
	if err != nil {
		log.Fatalln(err)
	}
}

func validate(c *gin.Context) {
	url := c.Param(paramUrl)
	if url == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	w, err := strconv.Atoi(c.Param(paramWidth))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	h, err := strconv.Atoi(c.Param(paramHeight))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Set(contextParams, cache.Params{
		Url:    fmt.Sprintf("http:/%s", url),
		Width:  w,
		Height: h,
	})
}

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
