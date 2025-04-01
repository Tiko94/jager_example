package main

import (
	"time"
)

var (
	nextCommentId  = 100
	nextBlogPostId = 100
)

func makeTime(year int, month time.Month, day int) *time.Time {
	result := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return &result
}
