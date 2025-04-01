package controller

import (
	"bytes"
	"context"
	"flag"
	"io"
	"net/http"
	"strconv"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	opentracing "github.com/opentracing/opentracing-go"
)

var cachServiceURL = flag.String("cacheservice", "http://cacheservice:5000", "Address of the caching service provider")

func getFromCache(ctx context.Context, key string) (io.ReadCloser, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "cache-processing")
	defer span.Finish()

	req, _ := http.NewRequest(http.MethodGet, *cachServiceURL+"/?key="+key, nil)
	util.InjectSpanToReq(span, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		span.SetTag("cache", false)
		return nil, false
	}

	span.SetTag("cache", true)
	return resp.Body, true
}

func saveToCache(ctx context.Context, key string, duration int64, data []byte) {
	span, _ := opentracing.StartSpanFromContext(ctx, "save-to-cache")
	defer span.Finish()

	req, _ := http.NewRequest(http.MethodPost, *cachServiceURL+"/?key="+key,
		bytes.NewBuffer(data))
	req.Header.Add("cache-control", "maxage="+strconv.FormatInt(duration, 10))

	util.InjectSpanToReq(span, req)
	http.DefaultClient.Do(req)
}

func invalidateCacheEntry(ctx context.Context, key string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "invalidate-cache-entry")
	defer span.Finish()

	req, _ := http.NewRequest(http.MethodGet, *cachServiceURL+"/invalidate?key="+key, nil)
	util.InjectSpanToReq(span, req)
	http.DefaultClient.Do(req)
}
