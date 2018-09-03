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
	)

//имя фамилия, отображаемое имя, аффка
//отдельная таблица связанная с юзером показывающая что и где он выбрал
//Модель вопросов Question{A,B}


type MyType struct {
	ID      uint   `gorm:"primary_key",json:"id"`
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}

var db, err = gorm.Open("sqlite3", "database.db")


func generateRefresh() string{

	signer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"name": "DICK",
		"iss":  "admin",
		"exp":  time.Now().Add(time.Minute * 20).Unix(),
	})

	tokenString, err := signer.SignedString([]byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"))

	fmt.Println(tokenString,err)
	return string(tokenString)
}


func Decode(tokenString string){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("b093be924f51ddfe2dcbd5eb69aa195b14dca0ad2325e9b3d56ded6c7c519e2c"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["name"], claims["exp"])
	} else {
		fmt.Println(err)
	}
}

func generateAccess() {

	signer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"name": "name",
		"iss":  "admin",
		"exp":  time.Now().Add(time.Minute * 20).Unix(),
	})

	tokenString, err := signer.SignedString([]byte(""))

	fmt.Println(tokenString,err)
}

func main() {

	generateAccess()
	Decode(generateRefresh())

	if err != nil {
		panic(err)
	}
	//close DB after Main function's end
	defer db.Close()
	migrate()
	router := mux.NewRouter()
	//Read does not need protection
	router.HandleFunc("/read", read).Methods("GET")
	router.Handle("/create", MiddlewareAuth(create)).Methods("POST")
	router.Handle("/update", MiddlewareAuth(update)).Methods("PUT")
	router.Handle("/delete", MiddlewareAuth(del)).Methods("DELETE")
	println("Server started on port 8080")
	http.ListenAndServe(":8080", Log(router))
	println("Exit")
}

func read(w http.ResponseWriter, r *http.Request) {
	var req []MyType


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
	db.Model(MyType{}).Where("id=?", req.ID).Update(req)
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
func deleteAfter(req *MyType, seconds int) {
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
			http.Redirect(w, r, "/error", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
		//Conditions after
	})
}

func checkToken(r *http.Request) bool {
	sec, _ := r.Cookie("secure")
	if sec.Value == "secure" {
		return true
	}
	return false
}

func GetDecodedJson(r *http.Request) (t MyType) {
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
	if (!db.HasTable(MyType{})) {
		println("migrated")
		db.AutoMigrate(MyType{})
		for i := 0; i < 10; i++ {
			seed := MyType{
				Name:    "random " + strconv.Itoa(i),
				Surname: "random " + strconv.Itoa(i),
			}
			db.Create(&seed)
		}
	}
}
