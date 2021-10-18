package main

import (
	"fmt"
	"go-rest-api/application/db"
	"go-rest-api/application/handlers"
	"go-rest-api/middleware"
	"io"
	"net/http"
	. "os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	file, err := OpenFile("todo.log", O_RDWR|O_CREATE|O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open file with error: " + err.Error())
	}

	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting ToDo API!!!")

	// database action
	db := db.DBInit()
	defer db.Close()
	// Untuk merestart database menjadi kosong
	// db.DropTableIfExists(&models.ToDoItem{}, &models.User{})
	// db.AutoMigrate(&models.ToDoItem{}, &models.User{})

	router := mux.NewRouter()
	// middleware
	router.Use(middleware.Auth(db))

	router.HandleFunc("/healthcheck", healthcheckHandler).Methods("GET")
	// user route
	router.HandleFunc("/api/v1/register", handlers.CreateUseHandler(db)).Methods("POST")
	router.HandleFunc("/api/v1/users", handlers.GetListUserhandler(db)).Methods("GET")

	// todo route
	router.HandleFunc("/api/v1/todo", handlers.CreateToDoHandler(db)).Methods("POST")
	router.HandleFunc("/api/v1/todos", handlers.GetListTodoHandler(db)).Methods("GET")
	router.HandleFunc("/api/v1/todo/{id}", handlers.GetTodoByIDHandler(db)).Methods("GET")
	router.HandleFunc("/api/v1/todo/{id}", handlers.DeleteTodoHandler(db)).Methods("DELETE")
	router.HandleFunc("/api/v1/todo/{id}", handlers.UpdateTodoHandler(db)).Methods("PUT")

	// Image Handler
	router.HandleFunc("/image/{imageName}", handlers.ShowImageHandler).Methods("GET")
	http.ListenAndServe(":8001", router)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	io.WriteString(w, `{"alive":true}`)
}
