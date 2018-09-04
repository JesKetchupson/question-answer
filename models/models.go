package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

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
type Answer struct {
	ID uint `gorm:"primary_key",json:"id"`
	User       *User `gorm:"unique;foreignkey:UserID"`
	UserID     uint `gorm:"type:INT UNSIGNED REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT",json:"user_id"`
	Objects    *[]Object `gorm:"foreignkey:QuestionID"`
	ObjectID   uint `gorm:"type:INT UNSIGNED REFERENCES objects(id) ON DELETE RESTRICT ON UPDATE RESTRICT",json:"object_id"`
	Questions  *[]Question `gorm:"foreignkey:QuestionID"`
	QuestionID uint `gorm:"type:INT UNSIGNED REFERENCES questions(id) ON DELETE RESTRICT ON UPDATE RESTRICT",json:"question_id"`
}

type Question struct {
	ID           uint   `gorm:"primary_key",json:"id"`
	FirstObject  *Object `gorm:"foreignkey:FirstObjectID"`
	FirstObjectID uint `gorm:"type:INT UNSIGNED REFERENCES objects(id) ON DELETE RESTRICT ON UPDATE RESTRICT",json:"first_object"`
	SecondObject *Object `gorm:"foreignkey:SecondObjectID"`
	SecondObjectID uint `gorm:"type:INT UNSIGNED REFERENCES objects(id) ON DELETE RESTRICT ON UPDATE RESTRICT",json:"second_object"`
}

type Object struct {
	ID         uint     `gorm:"primary_key",json:"id"`
	Category *Category `gorm:"foreignkey:CategoryID;association_foreignkey:ID"`
	CategoryID uint `gorm:"type:INT UNSIGNED REFERENCES categories(id) ON DELETE RESTRICT ON UPDATE RESTRICT",json:"category_id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
}

type Category struct {
	ID   uint   `gorm:"primary_key",json:"id"`
	Name string `json:"name"`
}

func (user *User) GenerateRefresh(email string, password string) {

	signer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":    email,
		"password": password,
		"iss":      "admin",
		"exp":      time.Now().Add(time.Minute * 20).Unix(),
	})

	tokenString, err := signer.SignedString([]byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"))
	println(tokenString, err)
	user.RefreshToken = tokenString
}
func (user *User) GenerateAccess(email string, password string) {
	signer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":    email,
		"password": password,
		"iss":      "admin",
		"exp":      time.Now().Add(time.Minute * 20).Unix(),
	})
	tokenString, err := signer.SignedString([]byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"))
	println(tokenString, err)
	user.AccessToken = tokenString

}
