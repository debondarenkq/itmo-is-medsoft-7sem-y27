package models

import "time"

type Staff struct {
	SecurityRevision int64      `json:"security_revision"`
	Revision         int64      `json:"identity_revision"`
	Role             string     `json:"role"`
	ID               string     `json:"id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Position         string     `json:"position"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}
