package helpers

import (
	. "awesomeProject/api/models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func Decode(tokenString string) User {
	token, errparse := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// []byte("my_secret_key")
		return []byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"), nil
	})

	if errparse != nil {
		panic(errparse)
	}
	var user User

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user = User{
			Email:    claims["email"].(string),
			Password: claims["password"].(string),
		}
	}
	return user
}

func GetDb() (*gorm.DB, error) {
	var db, err = gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	return db, err
}

/*
* func deleteAfter()
* Takes struct and int
* Delete record according to struct after N second
 */
func DeleteAfter(req *User, seconds int) {
	db, _ := GetDb()
	select {
	case <-time.After(time.Second * time.Duration(seconds)):
		db.Delete(&req)
	}
}
func GetDecodedJson(r *http.Request) (t User) {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&t)
	return
}
