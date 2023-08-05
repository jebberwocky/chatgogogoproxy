package chatHistory

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var tempCache *cache.Cache

func init() {
	tempCache = cache.New(5*time.Minute, 10*time.Minute)
}

func SetRule(pk string, rule map[string]interface{}) {
	_, found := tempCache.Get(pk)
	var messages []interface{}
	if !found {
		messages = append(messages, rule)
		tempCache.Set(pk, messages, cache.DefaultExpiration)
	}
}

func Set(pk string, message map[string]interface{}) []interface{} {
	foo, found := tempCache.Get(pk)
	var messages []interface{}
	if found {
		messages = foo.([]interface{})
	}
	messages = append(messages, message)
	tempCache.Set(pk, messages, cache.DefaultExpiration)
	return messages
}

func Get(pk string) []interface{} {
	foo, found := tempCache.Get(pk)
	if found {
		return foo.([]interface{})
	} else {
		return []interface{}{}
	}
}
