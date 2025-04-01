package entity

import "github.com/jinzhu/gorm"

type Author struct {
	gorm.Model
	FirstName string
	LastName  string
	Username  string
}
