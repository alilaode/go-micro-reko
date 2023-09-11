package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Hash      string    `json:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time `json:"modified_at" gorm:"<-:update"`
}

func (User) TableName() string {
	return "user"
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSession struct {
	JWTToken string `json:"jwt_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
