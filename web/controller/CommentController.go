package controller

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/web/model"
	opentracing "github.com/opentracing/opentracing-go"
)

type CommentController struct{}

func (c *CommentController) createComment(w http.ResponseWriter, r *http.Request) {
	rootSpan := util.GetSpanFromRPCReq(Tracer, r, "create-comment")
	defer rootSpan.Finish()

	matches := commentsPath.FindStringSubmatch(r.URL.Path)

	//no need to check for error since regex guarantees an integer value
	postID, _ := strconv.Atoi(matches[1])

	r.ParseForm()
	now := time.Now()
	comment := &entity.Comment{
		Subject:     r.FormValue("subject"),
		Body:        r.FormValue("body"),
		CreatedDate: now,
		PublishDate: now,
		PostID:      uint(postID),
		IsPublished: true,
	}

	comment, err := model.CreateComment(
		opentracing.ContextWithSpan(r.Context(), rootSpan),
		comment,
		postID,
	)

	if err != nil {
		rootSpan.SetTag("error", true)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("Location", "/posts/"+strconv.Itoa(postID))
	w.WriteHeader(http.StatusSeeOther)
	go invalidateCacheEntry(opentracing.ContextWithSpan(r.Context(), rootSpan), w.Header().Get("Location"))
}
