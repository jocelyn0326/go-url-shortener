package models

import (
	"time"
)

type UrlElement struct {
	LongUrl        string    `json:"longUrl" binding:"required"`
	ShortUrl       string    `json:"shortUrl"`
	CreateDateTime time.Time `json:"createdAt"`
}

type ErrorTemplate struct {
	Error string `json:"error" binding:"required"`
}
