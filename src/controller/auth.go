package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"backend-starter/src/models"
	"backend-starter/src/utils"
)

var users = make(map[string]string)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := models.Validate.Struct(user); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			fieldError := fmt.Sprintf("Поле %s: ", err.Field())
			switch err.Tag() {
			case "required":
				fieldError += "обязательно для заполнения"
			case "min":
				fieldError += fmt.Sprintf("должно содержать минимум %s символов", err.Param())
			case "email":
				fieldError += "должно быть корректным email адресом"
			default:
				fieldError += err.Tag()
			}
			validationErrors = append(validationErrors, fieldError)
		}
		http.Error(w, strings.Join(validationErrors, ", "), http.StatusBadRequest)
		return
	}

	if _, exists := users[user.Email]; exists {
		http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
		return
	}

	users[user.Email] = user.Password
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Регистрация успешна!"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if password, exists := users[user.Email]; !exists || password != user.Password {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})

	json.NewEncoder(w).Encode(models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, Message: "Вход успешен!"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Вы вышли из системы!"})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Необходима аутентификация", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Доступ разрешен!"))
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Необходима аутентификация", http.StatusUnauthorized)
		return
	}

	email := "test@gmail.com"

	accessToken, err := utils.GenerateAccessToken(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})

	json.NewEncoder(w).Encode(models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, Message: "Токен обновлен!"})
}
