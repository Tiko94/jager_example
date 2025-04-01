package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	dbConn  *gorm.DB
	initErr error
)

func createPost(ctx context.Context, blogPost entity.BlogPost) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql-create-post")
	defer span.Finish()

	span.SetTag("mysql", true)
	span.SetTag("sql.query", "INSERT INTO posts")

	return dbConn.Create(&blogPost).Error
}

func getPosts(ctx context.Context) ([]entity.BlogPost, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql-get-posts")
	defer span.Finish()

	span.SetTag("mysql", true)
	span.SetTag("sql.query", "SELECT * FROM posts")
	var blogPosts []entity.BlogPost
	db := dbConn.Find(&blogPosts)

	return blogPosts, db.Error
}

func getComments(ctx context.Context, postId uint) ([]entity.Comment, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql-get-comments")
	defer span.Finish()

	span.SetTag("mysql", true)
	span.SetTag("sql.query", "SELECT * FROM comments")
	var comments []entity.Comment
	db := dbConn.Find(&comments)

	return comments, db.Error
}

func getPost(ctx context.Context, id uint) (entity.BlogPost, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql-get-post")
	defer span.Finish()

	span.SetTag("mysql", true)
	span.SetTag("sql.query", "SELECT * FROM posts WHERE post_id = "+strconv.Itoa(int(id)))
	var blogPost entity.BlogPost
	dbConn.Preload("Comments").First(&blogPost, id)

	if blogPost.Subject == "" {
		return blogPost, errors.New("not found")
	}

	return blogPost, nil
}

func updatePost(ctx context.Context, id uint, post entity.BlogPost) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql-get-post")
	defer span.Finish()

	span.SetTag("mysql", true)
	span.SetTag("sql.query", "UPDATE posts WHERE post_id = "+strconv.Itoa(int(id)))
	var blogPost entity.BlogPost
	dbConn.First(&blogPost, id)

	if blogPost.Subject == "" {
		return errors.New("not found")
	}

	post.ID = blogPost.ID
	post.CreatedAt = blogPost.CreatedAt
	post.CreatedDate = blogPost.CreatedDate

	return dbConn.Save(post).Error
}

func newComment(ctx context.Context, id uint, comment entity.Comment) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "mysql-get-post")
	defer span.Finish()

	comment.PostID = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	span.SetTag("mysql", true)
	span.SetTag("sql.query", fmt.Sprintf("INSERT INTO comments (post_id = %v)", id))
	return dbConn.Create(&comment).Error
}
