package models

// User model for gorm
type User struct {
	ID        string
	CtzID     string
	Phone     string
	Prefix    string
	FirstName string
	LastName  string
	Plan      string
	ExamType  string
	Status    string
	Rank      string
	Confirmed bool
}
