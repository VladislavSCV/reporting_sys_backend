package models

import "time"

type Schedule struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"group_name"`
	SubjectName string    `json:"subject_name"`
	TeacherName string    `json:"teacher_name"`
	DayOfWeek   int       `json:"day_of_week"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     int       `json:"group_id"`
	SubjectID   int       `json:"subject_id"`
	TeacherID   int       `json:"teacher_id"`
}
