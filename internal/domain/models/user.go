package models

type User struct {
	ID          int64
	Email       string
	PassHash    []byte
	Username    string
	Sex         string
	Location    string
	DateOfBirth string
}
