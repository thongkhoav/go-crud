package models

import "gorm.io/gorm"

// post is a model for the post table
type Post struct {
	gorm.Model
	Title string 
	Body  string 
	UserID uint
	User   User `gorm:"foreignkey:UserID"`
}