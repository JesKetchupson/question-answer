package middleware

import (
	"net/http"
	"awesomeProject/api/helpers"
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
	sec, _ := r.Cookie("AccessToken")
	user := helpers.Decode(sec.Value)
	var db, _ = helpers.GetDb()
	isTrue := db.Where("email=?", user.Email).Where("password=?", user.Password)
	if isTrue.RecordNotFound() {
		return false
	}
	return true
}