package jwtmiddleware

import "time"

type User struct {
	Id       uint `json:"id"`
	UserType int  `json:"userType"`
}

type LoginResp struct {
	Expire   time.Time `json:"expire"`
	Token    string    `json:"token"`
	UserType int       `json:"userType"`
}

type UnauthorizedResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
