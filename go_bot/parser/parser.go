package parser

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"gopkg.in/mgo.v2/bson"
)

var CACHE *memcache.Client

func init() {
	CACHE = database.NewMemcache()
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func urlToKey(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

func IsCached(url string) bool {
	key := urlToKey(url)
	_, err := CACHE.Get(key)
	if err == nil {
		return true
	}
	return false
}

func CacheURL(url string) {
	key := urlToKey(url)
	CACHE.Set(&memcache.Item{
		Key:        key,
		Value:      []byte("true"),
		Expiration: (60 * 60) * 24,
	})
}

func DeferedCache(url string) {
	if err := recover(); err != nil {
		//os.Exit(1)
		Log.Error("%s", err)
	} else {
		CacheURL(url)
	}
}

func Titlelize(s string) string {
	return strings.Title(strings.ToLower(s))
}

func CreateMembermeship(DB database.MongoDB, member, organization models.Rel, source models.Source, role string, label string) {
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
			"role":         role,
			"label":        label,
		},
	}, &models.Membership{})
}

func LinkTo(resource string, id string) string {
	return "http://api.databr.io/v1/" + resource + "/" + id
}

func ToUtf8(iso8859_1 string) string {
	iso8859_1_buf := []byte(iso8859_1)
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
