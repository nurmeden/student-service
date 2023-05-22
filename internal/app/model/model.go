package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `json:"lastName"`
	Password  string             `json:"password"`
	Email     string             `json:"email"`
	Age       string             `json:"age"`
	Courses   []string           `bson:"courses" json:"courses"`
}

type SignInData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthToken struct {
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Course struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Students    []string `json:"students"`
}

type CourseResponse struct {
	CourseData Course `json:"data"`
}
