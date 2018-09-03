package controllers

import (
	"net/http"
	"awesomeProject/api/helpers"
	"encoding/json"
	. "awesomeProject/api/models"
)
var db,err = helpers.GetDb()

func Registration(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	db.NewRecord(req)
	db.Create(&req)

	req.AccessToken = helpers.GenerateAccess(req.Email, req.Password)

	req.RefreshToken = helpers.GenerateRefresh(req.Email, req.Password)

	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)

	isTrue := db.Where("email=?", req.Email).Where("password=?", req.Password)
	if !isTrue.RecordNotFound() {
		req.AccessToken = helpers.GenerateAccess(req.Email, req.Password)

		req.RefreshToken = helpers.GenerateRefresh(req.Email, req.Password)

		json.NewEncoder(w).Encode(req)
	}

	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Read(w http.ResponseWriter, r *http.Request) {
	var req []User
	db.Find(&req)

	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Create(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	db.NewRecord(req)
	db.Create(&req)
	json.NewEncoder(w).Encode(req)
}
func Update(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	db.Model(User{}).Where("id=?", req.ID).Update(req)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Del(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	go helpers.DeleteAfter(&req, 3)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
