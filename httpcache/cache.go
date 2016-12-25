package httpcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yookoala/weatherhk/ctxlog"
	rcache "gopkg.in/go-redis/cache.v5"
	redis "gopkg.in/redis.v5"
)

// Error represents error in httpcache
type Error int

const (
	// HeaderNotExists represents error if header field is empty
	// or does not exists
	HeaderNotExists Error = iota
)

func (err Error) Error() string {
	if err == HeaderNotExists {
		return "header field not exits"
	}
	return "unknown error"
}

// HKT stores *time.Location of Hong Kong
var HKT *time.Location

func init() {
	HKT, _ = time.LoadLocation("Asia/Hong_Kong")
}

const fmtRFC2612 = "Mon, 02 Jan 2006 15:04:05 GMT"

func parseRFC2612(str string) (t time.Time, err error) {
	if len(str) < 3 {
		err = fmt.Errorf("incorrect time string provided: %s", str)
		return
	}
	t, err = time.Parse(fmtRFC2612, str[:len(str)-3]+"GMT")
	if err != nil {
		// undo the string replacenemt for better error messages
		terr := err.(*time.ParseError)
		err = fmt.Errorf("cannot parse time %#v as RFC2612 (%#v)", str, terr.Layout)
	}
	return
}

var redisCache *rcache.Codec

func init() {
	redisURL, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Printf("invalid REDIS_URL: %s", err.Error())
		return
	}

	// assign global redisCache
	redisCache = &rcache.Codec{
		Redis: redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"default": redisURL.Addr,
			},
			Password: redisURL.Password,
		}),
		Marshal: func(v interface{}) ([]byte, error) {
			return json.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return json.Unmarshal(b, v)
		},
	}
}

// NewCache wraps an http.ResponswWriter with Cache
func NewCache(w http.ResponseWriter) *Cache {
	return &Cache{
		responseWriter: w,
		Created:        time.Now(),
		Status:         http.StatusOK,
		content:        bytes.NewBuffer(make([]byte, 0, 4096)), // pre=alloc 4096 bytes for buffer
	}
}

// Cache wraps an http.ResponseWriter and return
type Cache struct {
	responseWriter http.ResponseWriter
	content        *bytes.Buffer
	Created        time.Time
	Status         int
	CachedHeader   http.Header // header loaded from cache
	CachedContent  string
}

// Code returns the cached http status code
func (cache *Cache) Code() int {
	return cache.Status // TODO; ensure the default value is http.StatusOK
}

// TODO: add method to check Last-Modified date / Date of the cached response

// ParseTime parses the http header field with RFC2612 format
func (cache *Cache) ParseTime(name string) (parsed time.Time, err error) {
	var timeStr string
	if timeStr = cache.Header().Get(name); timeStr == "" {
		err = HeaderNotExists
		return
	}
	return parseRFC2612(timeStr)
}

// String return buffered content as string
func (cache *Cache) String() string {
	if cache.responseWriter != nil {
		return cache.content.String()
	}
	return cache.CachedContent
}

// Bytes return buffered content as string
func (cache *Cache) Bytes() []byte {
	if cache.responseWriter != nil {
		return cache.content.Bytes()
	}
	return []byte(cache.CachedContent)
}

// WriteTo writes the content of current cache to http.ResponseWriter
func (cache *Cache) WriteTo(w http.ResponseWriter) {
	header := cache.Header()
	for name := range header {
		for i := range header[name] {
			w.Header().Add(name, header[name][i])
		}
	}
	w.WriteHeader(cache.Code())
	w.Write(cache.Bytes())
}

// Header implements http.ResponseWriter
func (cache *Cache) Header() http.Header {
	if cache == nil {
		return nil
	}
	if cache.responseWriter != nil {
		return cache.responseWriter.Header()
	}
	return cache.CachedHeader
}

// Write implements http.ResponseWriter
func (cache *Cache) Write(p []byte) (int, error) {
	cache.content.Write(p) // omit write number and error in buffer write
	return cache.responseWriter.Write(p)
}

// WriteHeader implements http.ResponseWriter
func (cache *Cache) WriteHeader(code int) {
	cache.Status = code
	cache.responseWriter.WriteHeader(code)
}

func keyOf(r *http.Request) (key string, err error) {
	if r == nil {
		err = fmt.Errorf("request cannot be nil")
		return
	}
	if r.URL == nil {
		err = fmt.Errorf("r.URL cannot be nil")
		return
	}
	key = "page:/" + r.URL.Path
	return
}

// Load cache for a given http request
func Load(r *http.Request) (cache *Cache, err error) {
	key, err := keyOf(r)
	if err != nil {
		return
	}

	if redisCache == nil {
		return
	}

	cache = &Cache{}
	if err = redisCache.Get(key, cache); err != nil {
		cache = nil
		return
	}
	return
}

// Save cache for a given http request
func Save(r *http.Request, cache *Cache) (err error) {
	key, err := keyOf(r)
	if err != nil {
		return
	}

	if redisCache == nil {
		return
	}

	// ensure that the header is cached
	cache.CachedHeader = cache.Header()
	cache.CachedContent = cache.String()

	// store the httpcache item in memcached
	return redisCache.Set(&rcache.Item{
		Key:        key,
		Object:     cache,
		Expiration: 60 * time.Minute, // TODO: detect correct expiration time
	})
}

// Delete deletes cache of a given request
func Delete(r *http.Request) (err error) {
	key, err := keyOf(r)
	if err != nil {
		return
	}

	if redisCache == nil {
		return
	}

	return redisCache.Delete(key)
}

// Valid test if a cache has valid cache
func Valid(r *http.Request, cache *Cache) bool {
	var expires time.Time
	var err error
	infoLog, errorLog := ctxlog.GetLoggers(r)

	// TODO: might support max-age somehow?

	// parse grace expires override
	if expires, err = cache.ParseTime("X-Grace-Expires"); err == nil {
		if expires.After(time.Now()) {
			infoLog.Log("message", "cache graced")
			return true
		}
	} else if err != HeaderNotExists {
		errorLog.Log("message", fmt.Sprintf("error parsing X-Grace-Expires (%s)", err.Error()))
	}

	if expires, err = cache.ParseTime("Expires"); err == nil {
		if expires.After(time.Now()) {
			infoLog.Log("message", "cache not expired")
			return true
		}
	} else if err != HeaderNotExists {
		errorLog.Log("message", fmt.Sprintf("error parsing Expires (%s)", err.Error()))
	}
	return false // default treat as expired
}

// CacheHandler applies httpcache to the wrapped http.Handler
func CacheHandler(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		infoLog, errorLog := ctxlog.GetLoggers(r)

		// try to load cache for the request
		cache, err := Load(r)
		if err != nil {
			errorLog.Log("message", fmt.Sprintf("error loading cache: %s", err.Error()))
		}

		// if has cache, write to ResponseWriter and return early
		if Valid(r, cache) {
			infoLog.Log("message", "use cache")
			cache.WriteTo(w)
			return // early return
		}

		// refresh cache by running inner handler
		infoLog.Log("message", "no valid cache, trigger inner handler")

		cache = NewCache(w)
		inner.ServeHTTP(cache, r)
		go func() {
			err := Save(r, cache)
			if err != nil {
				errorLog.Log("message", fmt.Sprintf("error saving cache: %s", err.Error()))
			}
		}()
	})
}
