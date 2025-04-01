package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Subject     string
	Body        string
	PostID      uint
	Author      Author `gorm:"foreignkey:author_id"`
	CreatedDate time.Time
	PublishDate time.Time
	IsPublished bool
}
