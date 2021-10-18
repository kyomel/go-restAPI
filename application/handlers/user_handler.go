package handlers

import (
	"encoding/json"
	"go-rest-api/application/models"
	"net/http"

	"github.com/jinzhu/gorm"
)

func CreateUseHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		password := r.FormValue("password")

		newUser := &models.User{Name: name, Password: password}
		db.Create(&newUser)
		result := db.Last(&newUser)

		w.Header().Set("content-type", "appliacation/json")
		json.NewEncoder(w).Encode(result.Value)
	}
	return http.HandlerFunc(fn)
}

func GetListUserhandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		users := []models.User{}
		db.Find(&users)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
	return http.HandlerFunc(fn)
}
