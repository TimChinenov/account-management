package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Score    int    `json:"score"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
