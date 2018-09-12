package middleware

import (
	"awesomeProject/api/helpers"
	"net/http"
	)

func MiddlewareAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Conditions before
		if CheckToken(r) == false {
			//w.Write([]byte(strconv.Itoa(http.StatusForbidden)))
			http.Redirect(w, r, "/read", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
		//Conditions after
	})
}

func CheckToken(r *http.Request) bool {
	sec, err := r.Cookie("AccessToken")

	if err == nil{
	user,_ := helpers.Decode(sec.Value)
		var db, _ = helpers.GetDb()
		isTrue := db.Where("email=?", user.Email).Where("password=?", user.Password)
		if isTrue.RecordNotFound() {
			return false
		}
		return true
	}else {
		return false
	}

}
