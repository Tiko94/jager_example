package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	"github.com/jinzhu/gorm"
	opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)

const (
	envListenPort = "LISTEN_PORT"
	envDBHost     = "DB_HOST"
	envDBPort     = "DB_PORT"
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBName     = "DB_NAME"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	tracer, closer = util.InitJaeger("data-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	listen_port := getEnv(envListenPort, "4000")
	db_host := getEnv(envDBHost, "db")
	db_port := getEnv(envDBPort, "3306")
	db_user := getEnv(envDBUser, "root")
	db_password := getEnv(envDBPassword, "toor")
	db_name := getEnv(envDBName, "blog")

	connection_string := fmt.Sprintf("%s:%s@tcp(%s:%s)/%v?charset=utf8&parseTime=True&loc=Local",
		db_user, db_password, db_host, db_port, db_name)
	dbConn, initErr = gorm.Open("mysql", connection_string)
	if initErr != nil {
		log.Println(initErr)
		return
	}
	defer dbConn.Close()
	http.HandleFunc("/posts/", HandleRequest)
	http.HandleFunc("/posts", HandleRequest)
	http.HandleFunc("/healthcheck", HandleHealthcheck)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", listen_port), new(util.GzipHandler))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

var (
	postsPath    = regexp.MustCompile(`^/posts\?*`)
	postPath     = regexp.MustCompile(`^/posts/(\d+)`)
	commentsPath = regexp.MustCompile(`^/posts/(\d+)/comments`)
)

func HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "healthcheck")
	defer rootSpan.Finish()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"OK"}`))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}	
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		if commentsPath.MatchString(r.URL.Path) {
			HandleCreateComment(w, r)
		} else if postsPath.MatchString(r.URL.Path) {
			HandleCreatePost(w, r)
		}
	case http.MethodPut:
		if postPath.MatchString(r.URL.Path) {
			HandleUpdatePost(w, r)
		}
	case http.MethodGet:

		if commentsPath.MatchString(r.URL.Path) {
			HandleGetComments(w, r)
		} else if postPath.MatchString(r.URL.Path) {
			HandleGetPost(w, r)
		} else if postsPath.MatchString(r.URL.Path) {
			HandleGetPosts(w, r)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		
	}
}

func HandleGetComments(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "get-comments")
	defer rootSpan.Finish()

	matches := commentsPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postId, _ := strconv.Atoi(matches[1])

	comments, err := getComments(opentracing.ContextWithSpan(r.Context(), rootSpan), uint(postId))
	if err != nil {
		rootSpan.SetTag("error", true)
		rootSpan.LogFields(
			otlog.String("error-message", err.Error()),
		)

		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Unknown request"}`))
		if err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(comments); err != nil {
		log.Printf("Failed to encode comments: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "create-comment")
	defer rootSpan.Finish()

	matches := commentsPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postId, _ := strconv.Atoi(matches[1])

	dec := json.NewDecoder(r.Body)
	var comment entity.Comment
	err := dec.Decode(&comment)
	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		msg := `{"error":"` + err.Error() + `"}`
		_, writeErr := w.Write([]byte(msg))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
		return
	}

	if err := newComment(opentracing.ContextWithSpan(r.Context(), rootSpan), uint(postId), comment); err != nil {
		rootSpan.SetTag("error", true)
		rootSpan.LogFields(
			otlog.String("error-message", err.Error()),
		)

		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(comment); err != nil {
		log.Printf("Failed to encode comment: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "create-post")
	defer rootSpan.Finish()

	dec := json.NewDecoder(r.Body)
	var blogPost entity.BlogPost
	err := dec.Decode(&blogPost)
	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
		return
	}

	if err := createPost(opentracing.ContextWithSpan(r.Context(), rootSpan), blogPost); err != nil {
		rootSpan.SetTag("error", true)
		rootSpan.LogFields(
			otlog.String("error-message", err.Error()),
		)

		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(blogPost); err != nil {
		log.Printf("Failed to encode blog post: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "update-post")
	defer rootSpan.Finish()

	matches := postPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postId, _ := strconv.Atoi(matches[1])

	var blogPost entity.BlogPost
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&blogPost)
	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
		return
	}

	if err := updatePost(opentracing.ContextWithSpan(r.Context(), rootSpan), uint(postId), blogPost); err != nil {
		rootSpan.SetTag("error", true)
		rootSpan.LogFields(
			otlog.String("error-message", err.Error()),
		)

		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(&blogPost); err != nil {
		log.Printf("Failed to encode blog post: %v", err)
	}
}

func HandleGetPost(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "get-post")
	defer rootSpan.Finish()

	matches := postPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postId, _ := strconv.Atoi(matches[1])
	blogPost, err := getPost(opentracing.ContextWithSpan(r.Context(), rootSpan), uint(postId))
	if err != nil {
		rootSpan.SetTag("error", true)
		rootSpan.LogFields(
			otlog.String("error-message", err.Error()),
		)

		w.WriteHeader(http.StatusNotFound)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(blogPost); err != nil {
		log.Printf("Failed to encode blog post: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(tracer, r, "get-posts")
	defer rootSpan.Finish()

	blogPosts, err := getPosts(opentracing.ContextWithSpan(r.Context(), rootSpan))
	if err != nil {
		rootSpan.SetTag("error", true)
		rootSpan.LogFields(
			otlog.String("error-message", err.Error()),
		)

		w.WriteHeader(http.StatusNotFound)
		_, writeErr := w.Write([]byte(`{"error": "Unknown request"}`))
		if writeErr != nil {
			log.Printf("Failed to write error response: %v", writeErr)
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(blogPosts); err != nil {
		log.Printf("Failed to encode blog post: %v", err)
	}
}
// vim: tabstop=8 shiftwidth=8 expandtab!
