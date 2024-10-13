package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"backend-starter/src/controller"
	"backend-starter/src/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", controller.LogoutHandler).Methods("POST")
	r.HandleFunc("/refresh", controller.RefreshTokenHandler).Methods("POST")
	r.HandleFunc("/protected", middleware.TokenVerifyMiddleware(controller.ProtectedHandler)).Methods("GET")

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
