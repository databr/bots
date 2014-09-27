package parser

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/databr/api/models"
	"gopkg.in/mgo.v2/bson"
)

var CACHE *memcache.Client

func init() {
	memcacheURL := os.Getenv("MEMCACHE_URL")
	CACHE = memcache.New(memcacheURL)
	CACHE.Set(&memcache.Item{Key: "test", Value: []byte("tested")})
	_, err := CACHE.Get("test")
	if err != nil && err != memcache.ErrCacheMiss {
		panic(err)
	}
}

type Parser interface {
	Run(DB models.Database)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func urlToKey(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

func isCached(url string) bool {
	key := urlToKey(url)
	_, err := CACHE.Get(key)
	if err == nil {
		return true
	}
	return false
}

func cacheURL(url string) {
	key := urlToKey(url)
	CACHE.Set(&memcache.Item{
		Key:        key,
		Value:      []byte("true"),
		Expiration: (60 * 60) * 24,
	})
}

func deferedCache(url string) {
	if err := recover(); err != nil {
		// os.Exit(1)
		log.Error("%s", err)
	} else {
		cacheURL(url)
	}
}

func titlelize(s string) string {
	return strings.Title(strings.ToLower(s))
}

func createMembermeship(DB models.Database, member, organization models.Rel, source models.Source) {
	query := bson.M{
		"member.id":       member.Id,
		"organization.id": organization.Id,
	}
	DB.Upsert(query, bson.M{
		"$setOnInsert": bson.M{
			"createdat": time.Now(),
		},
		"$currentDate": bson.M{
			"updatedat": true,
		},
		"$set": bson.M{
			"member":       member,
			"organization": organization,
			"source":       source,
		},
	}, &models.Membership{})
}

func LinkTo(resource string, id string) string {
	return "http://api.databr.io/v1/" + resource + "/" + id
}
