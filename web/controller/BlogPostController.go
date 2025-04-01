package controller

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/web/model"
	"github.com/jinzhu/gorm"
	opentracing "github.com/opentracing/opentracing-go"
)

type BlogPostController struct {
	blogListTemplate *template.Template
	blogTemplate     *template.Template
}

func (c *BlogPostController) showBlogList(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(Tracer, r, "show-blog-list")
	defer rootSpan.Finish()

	w.Header().Add("Content-Type", "text/html")

	cacheKey := url.QueryEscape(r.URL.RequestURI())
	resp, ok := getFromCache(opentracing.ContextWithSpan(r.Context(), rootSpan), cacheKey)

	if ok {
		rootSpan.SetTag("cached", true)
		io.Copy(w, resp)
		resp.Close()
		return
	}

	posts, err := model.GetLastPosts(opentracing.ContextWithSpan(r.Context(), rootSpan), 3)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	titles, err := model.GetLastPostTitles(opentracing.ContextWithSpan(r.Context(), rootSpan), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	context := map[string]interface{}{
		"posts":  posts,
		"titles": titles,
	}

	rootSpan.LogEvent("sending data")
	buf := bytes.Buffer{}
	c.blogListTemplate.Execute(&buf, context)
	data := buf.Bytes()
	w.Write(data)

	go saveToCache(opentracing.ContextWithSpan(r.Context(), rootSpan), cacheKey, int64(24*time.Hour), data[:])
}

func (c *BlogPostController) showBlogPost(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(Tracer, r, "show-blog-post")
	defer rootSpan.Finish()

	w.Header().Add("Content-Type", "text/html")

	cacheKey := url.QueryEscape(r.URL.RequestURI())
	resp, ok := getFromCache(opentracing.ContextWithSpan(r.Context(), rootSpan), cacheKey)
	if ok {
		rootSpan.SetTag("cached", true)
		io.Copy(w, resp)
		resp.Close()
		return
	}

	matches := postPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postID, _ := strconv.Atoi(matches[1])

	post, err := model.GetPostById(opentracing.ContextWithSpan(r.Context(), rootSpan), postID)

	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	titles, err := model.GetLastPostTitles(opentracing.ContextWithSpan(r.Context(), rootSpan), 10)

	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	context := map[string]interface{}{
		"post":   post,
		"titles": titles,
	}

	buf := bytes.Buffer{}
	c.blogTemplate.Execute(&buf, context)
	data := buf.Bytes()
	w.Write(data)
	go saveToCache(opentracing.ContextWithSpan(r.Context(), rootSpan), cacheKey, int64(24*time.Hour), data[:])
}

func (c *BlogPostController) createBlogPost(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(Tracer, r, "create-blog-post")
	defer rootSpan.Finish()

	r.ParseForm()
	now := time.Now()
	post := &entity.BlogPost{
		Subject:     r.FormValue("subject"),
		Body:        r.FormValue("body"),
		Comments:    []entity.Comment{},
		CreatedDate: now,
		IsPublished: false,
	}

	post, err := model.CreateBlogPost(opentracing.ContextWithSpan(r.Context(), rootSpan), post)

	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("Location", "/posts/"+strconv.Itoa(int(post.ID)))
	w.WriteHeader(http.StatusSeeOther)
}

func (c *BlogPostController) updateBlogPost(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(Tracer, r, "update-blog-post")
	defer rootSpan.Finish()

	matches := postPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postID, _ := strconv.Atoi(matches[1])

	r.ParseForm()
	now := time.Now()
	post := &entity.BlogPost{
		Model: gorm.Model{
			ID: uint(postID),
		},
		Subject:     r.FormValue("subject"),
		Body:        r.FormValue("body"),
		Comments:    []entity.Comment{},
		CreatedDate: now,
		IsPublished: false,
	}

	post, err := model.UpdateBlogPost(opentracing.ContextWithSpan(r.Context(), rootSpan), post)

	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("Location", "/posts/"+strconv.Itoa(int(post.ID)))
	w.WriteHeader(http.StatusSeeOther)
	go invalidateCacheEntry(opentracing.ContextWithSpan(r.Context(), rootSpan), url.QueryEscape(r.RequestURI))
}
