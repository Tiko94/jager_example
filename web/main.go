package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/logservice/loghelper"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/web/controller"
	opentracing "github.com/opentracing/opentracing-go"
)

var loadbalancerURL = flag.String("loadbalancer", "http://loadbalancer:2001", "Address of the load balancer")

func main() {
	port := os.Getenv("PORT")
	controller.Tracer, controller.Closer = util.InitJaeger(fmt.Sprintf("frontend-%s", port))
	defer controller.Closer.Close()

	opentracing.SetGlobalTracer(controller.Tracer)

	flag.Parse()

	templateCache, _ := buildTemplateCache()
	controller.Setup(templateCache)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		rootSpan := util.GetSpanFromRPCReq(controller.Tracer, r, "healthcheck")
		defer rootSpan.Finish()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"OK"}`))
	})

	go http.ListenAndServe(fmt.Sprintf(":%v", port), new(util.GzipHandler))

	go func() {
		for range time.Tick(300 * time.Millisecond) {
			tc, isUpdated := buildTemplateCache()
			if isUpdated {
				controller.SetTemplateCache(tc)
			}
		}
	}()
	time.Sleep(1 * time.Second)
	go loghelper.WriteEntry(&entity.LogEntry{
		Level:     entity.LogLevelInfo,
		Timestamp: time.Now(),
		Source:    "app server",
		Message:   "Registering with load balancer",
	})
	http.Get(*loadbalancerURL + fmt.Sprintf("/register?port=%s", port))

	log.Println("Server started, press <ENTER> to exit")

	waitCh := make(chan struct{})
	<-waitCh

	go loghelper.WriteEntry(&entity.LogEntry{
		Level:     entity.LogLevelInfo,
		Timestamp: time.Now(),
		Source:    "app server",
		Message:   "Unregistering with load balancer",
	})
	http.Get(*loadbalancerURL + fmt.Sprintf("/unregister?port=%s", port))
}

var lastModTime time.Time = time.Unix(0, 0)

func buildTemplateCache() (*template.Template, bool) {
	needUpdate := false

	f, _ := os.Open("/go/static/templates")

	fileInfos, _ := f.Readdir(-1)
	fileNames := make([]string, len(fileInfos))
	for idx, fi := range fileInfos {
		if fi.ModTime().After(lastModTime) {
			lastModTime = fi.ModTime()
			needUpdate = true
		}
		fileNames[idx] = "/go/static/templates/" + fi.Name()
	}

	var tc *template.Template
	if needUpdate {
		log.Print("Template change detected, updating...")
		tc = template.Must(template.New("").ParseFiles(fileNames...))
		log.Println("template update complete")
	}
	return tc, needUpdate
}
