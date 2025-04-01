package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BlogPost struct {
	gorm.Model
	Subject     string
	Body        string
	Author      Author    `gorm:"foreignkey:author_id"`
	Comments    []Comment `gorm:"foreignkey:PostID"`
	CreatedDate time.Time
	PublishDate time.Time
	IsPublished bool
}

func (BlogPost) TableName() string {
	return "posts"
}
