package common

import "github.com/mitchellh/hashstructure/v2"

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
