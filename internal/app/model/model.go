package model

import "time"

type Student struct {
	ID        string   `bson:"_id,omitempty"`
	FirstName string   `bson:"firstName"`
	LastName  string   `bson:"lastName"`
	Password  string   `bson:"password"`
	Email     string   `json:"email"`
	Age       string   `bson:"age"`
	Courses   []string `json:"courses"`
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
