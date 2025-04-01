package main

import (
	"crypto/tls"
	"flag"
	"io"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	opentracing "github.com/opentracing/opentracing-go"
)

type webRequest struct {
	r      *http.Request
	w      http.ResponseWriter
	doneCh chan struct{}
}

var (
	tracer opentracing.Tracer
	closer io.Closer
)

var (
	requestCh    = make(chan *webRequest)
	registerCh   = make(chan string)
	unregisterCh = make(chan string)
	heartbeartCh = time.Tick(5 * time.Second)
)

var (
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
)

func init() {
	http.DefaultClient = &http.Client{Transport: &transport}
}

func main() {
	tracer, closer = util.InitJaeger("loadbalancer")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})
		requestCh <- &webRequest{r: r, w: w, doneCh: doneCh}
		<-doneCh
	})

	go processRequests()

	go http.ListenAndServe(":2000", nil)

	go http.ListenAndServe(":2001", new(appserverHandler))
	waitCh := make(chan struct{})
	<-waitCh
}

type appserverHandler struct{}

func (h *appserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	port := r.URL.Query().Get("port")
	switch r.URL.Path {
	case "/register":
		registerCh <- ip + ":" + port
	case "/unregister":
		unregisterCh <- ip + ":" + port
	case "/healthcheck":
		HandleHealthcheck(w, r)
	}
}

func HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "healthcheck")
	defer rootSpan.Finish()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"OK"}`))
}
