package main

import (
	"cache/cache"
	"fmt"
)

func main() {
	c := cache.NewCache(10)
	c.AddCache("key1", "test")
	cacheItem, b := c.GetCacheItem("key1")
	if b {
		fmt.Println(cacheItem)
	}
	fmt.Println(printDelItemCache(c.DeleteCacheItem("key2")))
	fmt.Println(printDelItemCache(c.DeleteCacheItem("key1")))
	c.Clean()
}
func printDelItemCache(b bool, err error) string {
	if err != nil {
		return fmt.Sprint(err)
	} else if b {
		return "item удален"
	}
	return "item не удалился"
}
