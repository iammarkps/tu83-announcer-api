package models

// User model for gorm
type User struct {
	ID        string
	CtzID     string
	FirstName string
	LastName  string
	Plan      string
	ExamType  string
	Status    bool
	Rank      uint16
}
