package cache

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/genboo/final_project/common"
)

const (
	extension = ".jpg"
)

type Params struct {
	Url    string
	Width  int
	Height int
}

type imageCache struct {
	mu       sync.Mutex
	cache    *Cache
	cacheDir string
}

var ImageCache *imageCache
var httpClient = client()

func Init(cacheDir string) {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	val, err := strconv.Atoi(os.Getenv("max_cache_capacity"))
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("capacity: %d\n", val)
	c := NewCache(val)
	c.OnEvicted = removeFile
	ImageCache = &imageCache{
		cache:    c,
		cacheDir: cacheDir,
	}

	files, err := os.ReadDir(cacheDir)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				key := strings.TrimRight(file.Name(), extension)
				ImageCache.cache.Set(key, filepath.Join(cacheDir, file.Name()))
				log.Println(key)
				log.Println(filepath.Join(cacheDir, file.Name()))
			}
		}
	}
}

func (ic *imageCache) GetImage(header http.Header, params Params) ([]byte, error) {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	hash, _ := common.Hash(fmt.Sprintf("%d-%d-%s", params.Width, params.Height, params.Url))
	cacheKey := strconv.FormatUint(hash, 10)
	if data, ok := ic.cache.Get(cacheKey); ok {
		log.Println("from cache")
		file, err := os.Open(data.(string))
		defer file.Close()
		if err != nil {
			return nil, err
		}
		return io.ReadAll(file)
	}
	log.Println("recreate")
	req, err := http.NewRequest(http.MethodGet, params.Url, nil)
	if err != nil {
		return nil, err
	}
	for k, h := range header {
		for _, v := range h {
			req.Header.Add(k, v)
		}
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	srcImage, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	var resizedImage image.Image
	if srcImage.Bounds().Dx() < srcImage.Bounds().Dy() || params.Width <= params.Height {
		resizedImage = imaging.Resize(srcImage, 0, params.Height, imaging.Lanczos)
	} else {
		resizedImage = imaging.Resize(srcImage, params.Width, 0, imaging.Lanczos)
	}
	croppedImage := imaging.Thumbnail(resizedImage, params.Width, params.Height, imaging.Lanczos)

	cacheFilename := filepath.Join(ic.cacheDir, cacheKey+extension)
	err = imaging.Save(croppedImage, cacheFilename)
	if err != nil {
		return nil, err
	}
	ic.cache.Set(cacheKey, cacheFilename)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, croppedImage, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func removeFile(_ Key, value interface{}) {
	err := os.Remove(value.(string))
	log.Printf("remove %s\f", value.(string))
	if err != nil {
		log.Println(err)
	}
}

func client() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100
	return &http.Client{Timeout: time.Second * 10, Transport: transport}
}
