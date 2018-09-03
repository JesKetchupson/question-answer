package models

type User struct {
	ID           uint   `gorm:"primary_key",json:"id"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Email        string `gorm:"not null;unique",json:"email,omitempty"`
	Password     string `gorm:"not null",json:"password,omitempty"`
	Status       uint   `json:"status,omitempty"`
	AccessToken  string `gorm:"-",json:"access_token,omitempty"`
	RefreshToken string `gorm:"-",json:"refresh_token,omitempty"`
}
