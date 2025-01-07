package models

import "time"

type Attendance struct {
	ID        int       `json:"id"`
	StudentID int       `json:"student_id"`
	SubjectID int       `json:"subject_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` // present, absent, excused
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AttendanceDetail struct {
	ID          int       `json:"id"`
	StudentID   int       `json:"student_id"`
	StudentName string    `json:"student_name"`
	SubjectName string    `json:"subject_name"`
	Date        time.Time `json:"date"` // Используем time.Time вместо string
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
