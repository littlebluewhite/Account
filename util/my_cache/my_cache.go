package my_cache

import "time"

type cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, value interface{}, d time.Duration)
}

func GetCacheMap[T any, U comparable](c cache, name string) map[U]T {
	var result map[U]T
	if x, found := c.Get(name); found {
		result = x.(map[U]T)
	} else {
		return make(map[U]T)
	}
	return result
}
