package data

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/util"
	opentracing "github.com/opentracing/opentracing-go"
)

type BlogPostRepository interface {
	GetRecentTitles(ctx context.Context, count int) ([]*BlogSummary, error)
	GetRecentPosts(ctx context.Context, count int) ([]*entity.BlogPost, error)
	CreatePost(ctx context.Context, post *entity.BlogPost) (*entity.BlogPost, error)
	UpdatePost(ctx context.Context, post *entity.BlogPost) (*entity.BlogPost, error)
	GetById(ctx context.Context, postId int) (*entity.BlogPost, error)
}

func NewBlogPostRepository() BlogPostRepository {
	return &blogPostRepository{}
}

type blogPostRepository struct{}

type sortByDate []entity.BlogPost

func (s sortByDate) Len() int {
	return len(s)
}

func (s sortByDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByDate) Less(i, j int) bool {
	return s[i].PublishDate.After(s[j].PublishDate)
}

func (r *blogPostRepository) GetRecentTitles(ctx context.Context, count int) ([]*BlogSummary, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get-recent-titles")
	defer span.Finish()

	req, _ := http.NewRequest(http.MethodGet, *dataServiceUrl+"/posts", nil)
	util.InjectSpanToReq(span, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []entity.BlogPost
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&posts)

	if err != nil {
		return nil, err
	}

	sort.Sort(sortByDate(posts))

	result := []*BlogSummary{}
	for i := 0; i < count && i < len(posts); i++ {
		result = append(result, &BlogSummary{
			ID:         int(posts[i].ID),
			Subject:    posts[i].Subject,
			AuthorName: posts[i].Author.Username,
		})
	}

	return result, nil
}

func (r *blogPostRepository) GetRecentPosts(ctx context.Context, count int) ([]*entity.BlogPost, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get-recent-posts")
	defer span.Finish()

	req, _ := http.NewRequest(http.MethodGet, *dataServiceUrl+"/posts", nil)
	util.InjectSpanToReq(span, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []entity.BlogPost
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&posts)

	if err != nil {
		return nil, err
	}
	sort.Sort(sortByDate(posts))

	result := []*entity.BlogPost{}
	for i := 0; i < count && i < len(posts); i++ {
		result = append(result, &posts[i])
	}

	return result, nil
}

func (r *blogPostRepository) CreatePost(ctx context.Context, post *entity.BlogPost) (*entity.BlogPost, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "create-post")
	defer span.Finish()

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)

	err := enc.Encode(post)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, *dataServiceUrl+"/posts", &buf)
	req.Header.Add("Content-Type", "application/json")
	util.InjectSpanToReq(span, req)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *blogPostRepository) UpdatePost(ctx context.Context, post *entity.BlogPost) (*entity.BlogPost, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "update-post")
	defer span.Finish()

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)

	err := enc.Encode(post)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, *dataServiceUrl+"/posts/"+
		strconv.Itoa(int(post.ID)), &buf)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	util.InjectSpanToReq(span, req)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *blogPostRepository) GetById(ctx context.Context, postId int) (*entity.BlogPost, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get-by-id")
	defer span.Finish()

	req, _ := http.NewRequest(http.MethodGet, *dataServiceUrl+"/posts/"+strconv.Itoa(int(postId)), nil)
	util.InjectSpanToReq(span, req)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var post *entity.BlogPost
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&post)

	if err != nil {
		return nil, err
	}

	return post, nil
}
