package handlers

import (
	"encoding/json"
	"fmt"
	"go-rest-api/application/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func CreateToDoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		description := r.FormValue("description")

		// Menambahkan gambar file di ToDoList
		fileUploaded, header, err := r.FormFile("image")
		if err != nil {
			log.Warn("Failed get image with error " + err.Error())
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"Error": "Failed get image"}`)
			return
		}
		defer fileUploaded.Close()
		dir, err := os.Getwd()
		if err != nil {
			log.Warn("Failed get working directory " + err.Error())
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"Error": "Failed get working directory"}`)
			return
		}
		t := time.Now()
		fileName := fmt.Sprintf("%s%s%s", "image", t.Format("20060102150405"), filepath.Ext(header.Filename))
		fileLocation := filepath.Join(dir, "images", fileName)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Warn("Failed open file with error " + err.Error())
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"Error": "Failed open file"}`)
			return
		}

		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileUploaded)
		if err != nil {
			log.Warn("Failed copy file with error " + err.Error())
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"Error": "Failed copy file"}`)
			return
		}

		fullPathName := fmt.Sprintf("%s/image/%s", r.Host, fileName)
		newTodo := &models.ToDoItem{Description: description, IsCompleted: false, ImageURL: fullPathName}
		db.Create(&newTodo)
		result := db.Last(&newTodo)
		log.WithFields(log.Fields{"Description": description, "ImageURL": fullPathName}).Info("Success Add new Todo Item")

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(result.Value)
	}
	return http.HandlerFunc(fn)
}

func GetListTodoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		todoList := []models.ToDoItem{}

		db.Find(&todoList)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(todoList)
	}
	return http.HandlerFunc(fn)
}

func GetTodoByIDHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var todo models.ToDoItem
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		db.Where("id = ?", id).First(&todo)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
	return http.HandlerFunc(fn)
}

func UpdateTodoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var todo models.ToDoItem
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		description := r.FormValue("description")
		isCompleted, _ := strconv.ParseBool(r.FormValue("iscompleted"))

		db.Where("id = ?", id).First(&todo)
		todo.Description = description
		todo.IsCompleted = isCompleted

		db.Save(todo)
		log.WithFields(log.Fields{"ID": id, "Description": description, "IsCompleted": isCompleted}).Info("Success Updating ToDo Item")
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
	return http.HandlerFunc(fn)
}

func DeleteTodoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var todo models.ToDoItem
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		db.Where("id = ?", id).First(&todo)

		db.Delete(todo)
		log.WithFields(log.Fields{"ID": id}).Info("Success Delete ToDo Item")
		w.Header().Set("content-type", "json.NewEncoder(w).Encode(todo)application/json")
		io.WriteString(w, `{"Success": true}`)
	}
	return http.HandlerFunc(fn)
}
