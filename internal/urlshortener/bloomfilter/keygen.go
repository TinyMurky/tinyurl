package bloomfilter

// genBase62IDKey create key to store base62ID to bloom filter
func genBase62IDKey() string {
	key := "urlshortener:base62ID"
	return key
}

// genLongURLKey create key to store  to bloom filter
// func genLongURLKey() string {
// 	key := "urlshortener:longURL"
// 	return key
// }
