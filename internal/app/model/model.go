package model

import "time"

type Student struct {
	ID        string   `bson:"_id,omitempty"`
	FirstName string   `bson:"firstName"`
	LastName  string   `bson:"lastName"`
	Password  string   `bson:"password"`
	Email     string   `json:"email"`
	Age       int      `bson:"age"`
	Courses   []string `json:"courses"`
}

type SignInData struct {
	UserID   string `json:"userId" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthToken struct {
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type StudentInput struct {
	FirstName string   `json:"first_name" binding:"required"`
	LastName  string   `json:"last_name" binding:"required"`
	Email     string   `json:"email" binding:"required"`
	Password  string   `json:"password" binding:"required"`
	Courses   []string `json:"courses"`
}

type StudentUpdateInput struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Courses   []string `json:"courses"`
}
