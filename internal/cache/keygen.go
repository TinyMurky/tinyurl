package cache

import (
	"fmt"

	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
)

// genURLKey create key to store url in cache
func genURLKey(u model.URL) string {
	base62ID := u.GetIDBase62()
	key := fmt.Sprintf("urlshortener:url:base64ID:%s", base62ID)
	return key
}
