package cache

import (
	"bytes"
	"errors"
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
	"github.com/genboo/final_project/hash"
)

const (
	extension = ".jpg"
)

var (
	ErrNotFound = errors.New("изображение не найдено")
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

func Init(cacheDir string, capacity int) {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
	c := NewCache(capacity)
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
			}
		}
	}
}

func (ic *imageCache) GetImage(header http.Header, params Params) ([]byte, error) {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	h, _ := hash.Make(fmt.Sprintf("%d-%d-%s", params.Width, params.Height, params.Url))
	cacheKey := strconv.FormatUint(h, 10)
	if data, ok := ic.cache.Get(cacheKey); ok {
		file, err := os.Open(data.(string))
		defer file.Close()
		if err != nil {
			return nil, err
		}
		return io.ReadAll(file)
	}

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

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.New(resp.Status)
	}

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

func removeFile(value interface{}) {
	err := os.Remove(value.(string))
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
