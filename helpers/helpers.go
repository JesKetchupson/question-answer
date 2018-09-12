package helpers

import (
	. "awesomeProject/api/models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"os"
	)

func Decode(tokenString string) (User, error) {
	token, errparse := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// []byte("my_secret_key")
		return []byte(os.Getenv("Secret")), nil
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
	return user, nil
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

type Configuration struct {
	Port              int
	Secret            string
}

func InitEnvVal(filename string)  {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fi,_ := os.Stat(filename)
	b := make([]byte,fi.Size())
	file.Read(b)

	c:=make(map[string]string)

	err = json.Unmarshal(b,&c)
	for key,val:=range c{
	os.Setenv(key,val)
	}

	if err != nil {
		panic(err)
	}
}
