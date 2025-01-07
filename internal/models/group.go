package models

import "time"

type Group struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	StudentCount int       `json:"student_count"`
}

type GroupDetail struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Students  []User     `json:"students"`
	Schedule  []Schedule `json:"schedule"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
