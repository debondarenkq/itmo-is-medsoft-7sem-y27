package models

import (
	"errors"
	"time"
)

var (
	ErrUnauthorized = errors.New("Неверный логин, пароль или одноразовый код")
	ErrForbidden    = errors.New("Недостаточно прав")
)

type Account struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
	StaffID      *string   `json:"staff_id,omitempty"`
	IsAdmin      bool      `json:"-"`
	TOTPSecret   string    `json:"-"`
	LastTOTPStep int64     `json:"-"`
	Role         string    `json:"role"`
	StaffName    string    `json:"staff_name"`
}

type AuditEvent struct {
	EventID   string    `json:"event_id"`
	EventTime time.Time `json:"event_time"`
	Actor     string    `json:"actor"`
	Message   string    `json:"message"`
	EventType string    `json:"event_type"`
	ActorName string    `json:"actor_name"`
}

type Claims struct {
	Subject  string `json:"sub"`
	StaffID  string `json:"staff_id,omitempty"`
	Role     string `json:"role"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Expires  int64  `json:"exp"`
	Issued   int64  `json:"iat"`
	ID       string `json:"jti"`
}

type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	AccountID    string    `json:"account_id"`
	StaffID      *string   `json:"staff_id,omitempty"`
	Role         string    `json:"role"`
}

type AccountChanges struct {
	Username     *string
	IsActive     *bool
	Role         *string
	PasswordHash string
}

type AccountFilter struct {
	Limit    int
	Offset   int
	Search   string
	IsActive *bool
}

type AuditFilter struct {
	Limit     int
	Offset    int
	EventType string
	From      *time.Time
	To        *time.Time
}
