package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/hashstructure/v2"
)

const (
	defaultTimeout = 10 * time.Second // На запрос максимум 10 секунд
	paramUrl       = "any"
	paramWidth     = "width"
	paramHeight    = "height"
	contextParams  = "params"
	previewPath    = "preview"
)

type api struct {
}

type Params struct {
	Url    string
	Width  int
	Height int
}

var httpClient = client()

var API = &api{}

func (v *api) Init(server *gin.Engine) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic occurred: %s\n", err)
		}
	}()
	server.GET("/fill/:width/:height/*any", v.validate, v.preview)
	if _, err := os.Stat(previewPath); os.IsNotExist(err) {
		if err = os.Mkdir(previewPath, 0o755); err != nil {
			log.Fatalln(err)
		}
	}
}

func (v *api) preview(c *gin.Context) {
	params := c.MustGet(contextParams).(*Params)
	// посмотреть в кэше
	hash, err := Hash(params)
	if err != nil {
		log.Fatalln(err)
	}
	key := strconv.FormatUint(hash, 10)
	if buf, err := Cache.Get(key); err == nil {
		_, err = c.Writer.Write(buf.Bytes())
	}
	// если нет, то загрузить и сохранить в кэше
	data, err := loadImage(c.Request.Header, params.Url)
	if err != nil {
		e := &Error{}
		if errors.As(err, &e) {
			c.AbortWithStatus(http.StatusBadGateway)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	err = resizePhoto(data, params.Width, params.Height)
	if err != nil {
		log.Fatalln(err)
	}
	Cache.Put()
	_, err = c.Writer.Write(data.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
}

func (v *api) validate(c *gin.Context) {
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
	c.Set(contextParams, &Params{
		Url:    url,
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

func Hash(v interface{}) (uint64, error) {
	return hashstructure.Hash(v, hashstructure.FormatV2, nil)
}
