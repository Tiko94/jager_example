package model

import (
	"context"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/web/data"
)

var (
	blogPostRepository data.BlogPostRepository = data.NewBlogPostRepository()
	commentRepository  data.CommentRepository  = data.NewCommentRepository()
)

func GetLastPosts(ctx context.Context, count int) ([]*entity.BlogPost, error) {
	return blogPostRepository.GetRecentPosts(ctx, count)
}

func GetLastPostTitles(ctx context.Context, count int) ([]*data.BlogSummary, error) {
	return blogPostRepository.GetRecentTitles(ctx, count)
}

func CreateBlogPost(ctx context.Context, post *entity.BlogPost) (*entity.BlogPost, error) {
	return blogPostRepository.CreatePost(ctx, post)
}

func UpdateBlogPost(ctx context.Context, post *entity.BlogPost) (*entity.BlogPost, error) {
	return blogPostRepository.UpdatePost(ctx, post)
}

func GetPostById(ctx context.Context, postId int) (*entity.BlogPost, error) {
	return blogPostRepository.GetById(ctx, postId)
}

func CreateComment(ctx context.Context, comment *entity.Comment, postId int) (*entity.Comment, error) {
	return commentRepository.CreateComment(ctx, comment, postId)
}
