package common

import "github.com/mitchellh/hashstructure/v2"

const (
	DefaultCacheDir = "image_cache"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func Hash(v interface{}) (uint64, error) {
	// хэширование ключа
	return hashstructure.Hash(v, hashstructure.FormatV2, nil)
}
