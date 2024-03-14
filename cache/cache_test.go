package cache

import (
	image2 "image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("image exist in cache", func(t *testing.T) {
		Init("image_cache_test", 5)
		key := "test"
		ImageCache.cache.Set(key, "")
		val, ok := ImageCache.cache.Get(key)
		require.True(t, ok)
		require.NotNil(t, val)
	})
	t.Run("image not found", func(t *testing.T) {
		Init("image_cache_test", 5)
		image, err := ImageCache.GetImage(nil, Params{
			Url:    "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x5041.jpg",
			Width:  100,
			Height: 100,
		})
		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, image)
	})
	t.Run("remote server error", func(t *testing.T) {
		Init("image_cache_test", 5)
		image, err := ImageCache.GetImage(nil, Params{
			Url:    "https://test1222.com/",
			Width:  100,
			Height: 100,
		})
		require.NotNil(t, err)
		require.Nil(t, image)
	})
	t.Run("not image", func(t *testing.T) {
		Init("image_cache_test", 5)
		image, err := ImageCache.GetImage(nil, Params{
			Url:    "https://raw.githubusercontent.com/OtusGolang/final_project/master/03-image-previewer.md",
			Width:  100,
			Height: 100,
		})
		require.NotNil(t, err)
		require.ErrorIs(t, err, image2.ErrFormat)
		require.Nil(t, image)
	})
	t.Run("image found", func(t *testing.T) {
		Init("image_cache_test", 5)
		image, err := ImageCache.GetImage(nil, Params{
			Url:    "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
			Width:  100,
			Height: 100,
		})
		require.Nil(t, err)
		require.NotNil(t, image)
		_ = os.RemoveAll("image_cache_test")
	})
}
