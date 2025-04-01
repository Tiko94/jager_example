package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	opentracing "github.com/opentracing/opentracing-go"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment, postId int) (*entity.Comment, error)
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

type commentRepository struct{}

func (r *commentRepository) CreateComment(ctx context.Context, comment *entity.Comment, postId int) (*entity.Comment, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "create-comment")
	defer span.Finish()

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.Encode(comment)

	req, _ := http.NewRequest(http.MethodPost, *dataServiceUrl+"/posts/"+strconv.Itoa(postId)+"/comments", &buf)
	util.InjectSpanToReq(span, req)

	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {

		errorText, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(errorText))
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&comment)

	if err != nil {
		return nil, err
	}
	return comment, nil
}
