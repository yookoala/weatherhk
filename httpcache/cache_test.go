package httpcache_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/yookoala/weatherhk/ctxlog"
	"github.com/yookoala/weatherhk/httpcache"
)

const fmtRFC2612 = "Mon, 02 Jan 2006 15:04:05 GMT"

func rfc2616(t time.Time) string {
	// RFC2616: Tue, 15 Nov 1994 12:45:26 GMT
	// RFC1123: Mon, 02 Jan 2006 15:04:05 MST
	return t.In(time.UTC).Format(fmtRFC2612)
}

func TestCache_implementsResponseWriter(t *testing.T) {
	var cache http.ResponseWriter = &httpcache.Cache{}
	_ = cache // just to verify httpcache.Cache implements http.ResponseWriter
}

func TestCache(t *testing.T) {

	// handler to be wrapped
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Custom message")
		w.Header().Add("X-Special-Header", "some special header message")
		w.Header().Add("Expires", rfc2616(time.Now().Add(-60*time.Second))) // should be rfc2616 format
		w.WriteHeader(http.StatusPartialContent)
	})

	w := httpcache.NewCache(httptest.NewRecorder())
	handler.ServeHTTP(w, nil)

	if want, have := "Custom message", w.String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "Custom message", string(w.Bytes()); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "some special header message", w.Header().Get("X-Special-Header"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := http.StatusPartialContent, w.Code(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := true, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestCache_Expired(t *testing.T) {
	w := httpcache.NewCache(httptest.NewRecorder())

	// manually update "Expires" header and check again
	w.Header().Set("Expires", rfc2616(time.Now().Add(60*time.Second))) // should be rfc2616 format
	if want, have := false, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	w.Header().Set("Expires", rfc2616(time.Now().Add(-60*time.Second))) // should be rfc2616 format
	if want, have := true, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// way to pass grace time instruction without
	// polluting "Expires" header (custom "X-Grace-Expires" header)
	w.Header().Set("Expires", rfc2616(time.Now().Add(-60*time.Second)))        // should be rfc2616 format
	w.Header().Set("X-Grace-Expires", rfc2616(time.Now().Add(60*time.Second))) // should be rfc2616 format
	if want, have := false, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	w.Header().Del("X-Grace-Expires")
	w.Header().Set("Expires", "some non-sense")
	if want, have := true, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	w.Header().Del("Expires")
	w.Header().Set("X-Grace-Expires", "some non-sense")
	if want, have := true, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	w.Header().Del("Expires")
	w.Header().Del("X-Grace-Expires")
	if want, have := true, w.Expired(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestCache_WriteTo(t *testing.T) {
	w1 := httpcache.NewCache(httptest.NewRecorder())
	w1.Header().Add("X-Custom-Header", "value 1")
	w1.Header().Add("X-Custom-Header", "value 2")
	w1.Header().Add("X-Custom-Header", "value 3")
	w1.WriteHeader(http.StatusBadGateway)
	w1.Write([]byte("Hello content"))

	w2 := httptest.NewRecorder()
	w1.WriteTo(w2)
	if want, have := len(w2.Header()), len(w1.Header()); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}
	if want, have := len(w2.Header()["X-Custom-Header"]), len(w1.Header()["X-Custom-Header"]); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := w2.Header()["X-Custom-Header"][0], w1.Header()["X-Custom-Header"][0]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := w2.Header()["X-Custom-Header"][1], w1.Header()["X-Custom-Header"][1]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := w2.Header()["X-Custom-Header"][2], w1.Header()["X-Custom-Header"][2]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestLoad(t *testing.T) {
	if url := os.Getenv("REDIS_URL"); url == "" {
		t.Skip("REDIS_URL not set, test skipped")
	}

	_, err := httpcache.Load(nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	r, _ := http.NewRequest("GET", "", nil)
	r.URL = nil
	_, err = httpcache.Load(r)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCacheHandler(t *testing.T) {

	if url := os.Getenv("REDIS_URL"); url == "" {
		t.Skip("REDIS_URL not set, test skipped")
	}

	// handler to be wrapped
	calledHandler := 0
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Custom message")
		w.Header().Add("X-Special-Header", "some special header message")
		w.Header().Add("Expires", rfc2616(time.Now().Add(60*time.Second)))
		w.WriteHeader(http.StatusPartialContent)
		calledHandler += 1
	})

	// serve request for handler
	handler = httpcache.CacheHandler(handler)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/some/path.html", nil)
	handler.ServeHTTP(w, r)

	if want, have := 1, calledHandler; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// wait a bit for memcached to handle
	time.Sleep(10 * time.Millisecond)

	// check if the cache actually works
	handler.ServeHTTP(w, r)
	if want, have := 1, calledHandler; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// examine the cached content
	cache, err := httpcache.Load(r)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if cache == nil {
		t.Error("cache is nil")
		return
	}
	defer httpcache.Delete(r)

	if want, have := "Custom message", cache.String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "Custom message", string(cache.Bytes()); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "some special header message", cache.Header().Get("X-Special-Header"); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := http.StatusPartialContent, cache.Code(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

type testHandler struct {
	message string
	expires time.Duration
	wait    time.Duration
	called  int
}

func (handler *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.called += 1
	time.Sleep(handler.wait)
	go func() {
		fmt.Fprint(w, handler.message)
		w.Header().Add("X-Special-Header", "some special header message")
		w.Header().Add("Expires", rfc2616(time.Now().Add(handler.expires)))
		w.WriteHeader(http.StatusPartialContent)
	}()
}

func benchmarkCacheHandler(b *testing.B, handler http.Handler) {

	// request once to activate the cache
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/some/path1.html", nil)
	handler.ServeHTTP(w, r)
	defer httpcache.Delete(r)

	// wait a bit for cache to handle
	time.Sleep(500 * time.Millisecond)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
	}
	b.StopTimer()
}

func BenchmarkCacheHandler_dry(b *testing.B) {

	b.StopTimer()

	if url := os.Getenv("REDIS_URL"); url == "" {
		b.Skip("REDIS_URL not set, benchmark skipped")
	}

	b.SetParallelism(1)

	// handler to be wrapped
	inner := &testHandler{
		message: "custom message",
		wait:    10 * time.Millisecond,
		expires: 24 * time.Hour,
	}

	benchmarkCacheHandler(b, inner)
}

func BenchmarkCacheHandler(b *testing.B) {

	b.StopTimer()

	if url := os.Getenv("REDIS_URL"); url == "" {
		b.Skip("REDIS_URL not set, benchmark skipped")
	}

	// discard log output
	ctxlog.SetLogOutput(ioutil.Discard)
	defer ctxlog.SetLogOutput(os.Stdout)

	// handler to be wrapped
	inner := &testHandler{
		message: "custom message",
		wait:    10 * time.Millisecond,
		expires: 24 * time.Hour,
	}

	benchmarkCacheHandler(b, httpcache.CacheHandler(inner))

	if want, have := 1, inner.called; want != have {
		b.Logf("expected %#v, got %#v", want, have)
	}
}
