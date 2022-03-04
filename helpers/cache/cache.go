package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	c *cache.Cache
)

func InitCache() {
	c = cache.New(cache.NoExpiration, 5*time.Second)
}

func Cache() *cache.Cache {
	return c
}

func CreateKeyLockPayment(listenerId string, date string, bookingTime string) (key string) {
	key = listenerId + "-" + date + ":" + bookingTime
	return key
}
