package main

import (
	"net/http"
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"strconv"
	"time"
	"log"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/gorilla/handlers"
	"os"
)

//имя фамилия, отображаемое имя, аффка
//отдельная таблица связанная с юзером показывающая что и где он выбрал
//Модель вопросов Question{A,B}

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

var db, err = gorm.Open("sqlite3", "database.db")

func main() {

	if err != nil {
		panic(err)
	}
	//close DB after Main function's end
	defer db.Close()

	migrate()

	router := mux.NewRouter()

	router.HandleFunc("/registration", registration).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")

	//Read does not need protection
	router.HandleFunc("/read", read).Methods("GET")
	router.Handle("/create", MiddlewareAuth(create)).Methods("POST")

	router.Handle("/update", MiddlewareAuth(update)).Methods("PUT")

	router.Handle("/delete", MiddlewareAuth(del)).Methods("DELETE")
	println("Server started on port 8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router))
	println("Exit")
}

func registration(w http.ResponseWriter, r *http.Request) {
	req := GetDecodedJson(r)
	db.NewRecord(req)
	db.Create(&req)

	req.AccessToken = generateAccess(req.Email, req.Password)

	req.RefreshToken = generateRefresh(req.Email, req.Password)

	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	req := GetDecodedJson(r)

	isTrue := db.Where("email=?", req.Email).Where("password=?", req.Password)
	if !isTrue.RecordNotFound() {
		req.AccessToken = generateAccess(req.Email, req.Password)

		req.RefreshToken = generateRefresh(req.Email, req.Password)

		json.NewEncoder(w).Encode(req)
	}

	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}

func read(w http.ResponseWriter, r *http.Request) {
	var req []User
	db.Find(&req)

	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	req := GetDecodedJson(r)
	db.NewRecord(req)
	db.Create(&req)
	json.NewEncoder(w).Encode(req)
}
func update(w http.ResponseWriter, r *http.Request) {
	req := GetDecodedJson(r)
	fmt.Println(req)
	db.Model(User{}).Where("id=?", req.ID).Update(req)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func del(w http.ResponseWriter, r *http.Request) {
	req := GetDecodedJson(r)
	go deleteAfter(&req, 3)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}

/*
* func deleteAfter()
* Takes struct and int
* Delete record according to struct after N second
*/
func deleteAfter(req *User, seconds int) {
	select {
	case <-time.After(time.Second * time.Duration(seconds)):
		db.Delete(&req)
	}
}

func MiddlewareAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Conditions before
		if checkToken(r) == false {
			//w.Write([]byte(strconv.Itoa(http.StatusForbidden)))
			http.Redirect(w, r, "/read", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
		//Conditions after
	})
}

func checkToken(r *http.Request) bool {
	sec, _ := r.Cookie("AccessToken")
	user := Decode(sec.Value)

	isTrue := db.Where("email=?", user.Email).Where("password=?", user.Password)

	if isTrue.RecordNotFound() {
		return false
	}
	return true
}

func GetDecodedJson(r *http.Request) (t User) {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&t)
	return
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func migrate() {
	db.DropTable(User{})
	if (!db.HasTable(User{})) {
		println("migrated")
		db.AutoMigrate(User{})
		for i := 1; i < 11; i++ {
			seed := User{
				Email:    "email" + strconv.Itoa(i),
				Password: "pass" + strconv.Itoa(i),
			}
			db.Create(&seed)
		}
	}
}
func generateRefresh(email string, password string) string {

	signer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":    email,
		"password": password,
		"iss":      "admin",
		"exp":      time.Now().Add(time.Minute * 20).Unix(),
	})

	tokenString, err := signer.SignedString([]byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"))

	fmt.Println(tokenString, err)
	return string(tokenString)
}
func generateAccess(email string, password string) string {
	signer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":    email,
		"password": password,
		"iss":      "admin",
		"exp":      time.Now().Add(time.Minute * 20).Unix(),
	})
	tokenString, err := signer.SignedString([]byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"))

	fmt.Println(tokenString, err)
	return string(tokenString)
}

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
