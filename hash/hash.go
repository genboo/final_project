package hash

import "github.com/mitchellh/hashstructure/v2"

func Make(v interface{}) (uint64, error) {
	// хэширование ключа
	return hashstructure.Hash(v, hashstructure.FormatV2, nil)
}
