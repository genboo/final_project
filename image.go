package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/nfnt/resize"
)

func loadImage(header http.Header, url string) (*bytes.Buffer, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// передача заголовков в запрос
	for k, v := range header {
		for _, i := range v {
			req.Header.Add(k, i)
		}
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		return nil, &Error{
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func resizePhoto(buf *bytes.Buffer, width, height int) error {
	img, _, err := image.Decode(buf)
	if err != nil {
		return err
	}
	if img.Bounds().Dx() < img.Bounds().Dy() || width <= height {
		img = resize.Resize(0, uint(height), img, resize.Lanczos3)
	} else {
		img = resize.Resize(uint(width), 0, img, resize.Lanczos3)
	}

	startX := (img.Bounds().Dx() - width) / 2
	endX := img.Bounds().Dx() - startX
	startY := (img.Bounds().Dy() - height) / 2
	endY := img.Bounds().Dy() - startY
	img = img.(interface {
		SubImage(image.Rectangle) image.Image
	}).SubImage(image.Rect(startX, startY, endX, endY))
	buf.Reset()
	if err = jpeg.Encode(buf, img, nil); err != nil {
		return err
	}
	return nil
}

func client() *http.Client {
	var timeout = defaultTimeout
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100
	return &http.Client{Timeout: timeout, Transport: transport}
}
