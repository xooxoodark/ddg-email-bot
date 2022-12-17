package main

import (
	"time"
)

type Token struct {
	ID        uint
	TID       int64
	Token     string
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WaitOTP struct {
	ID        uint
	TID       int64
	UserName  string
	Expire    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
