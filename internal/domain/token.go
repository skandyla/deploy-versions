package domain

import (
	"errors"
	"time"
)

var ErrRefreshTokenExpired = errors.New("refresh token expired")

type RefreshSession struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
}
