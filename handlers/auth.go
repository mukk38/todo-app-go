package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-api/models"

	"github.com/dgrijalva/jwt-go"
)

var users []models.User
var secretKey = []byte("supersecretkey") // Bu daha güvenli bir şekilde saklanmalı.

// JWT token oluşturma
func GenerateToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Kullanıcı kaydı
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Şifre kontrolü (şifreyi hash'lemek daha güvenli olur)
	users = append(users, user)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

// Kullanıcı girişi ve JWT token oluşturma
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	for _, u := range users {
		if u.Username == user.Username && u.Password == user.Password {
			token, err := GenerateToken(user.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Could not generate token"})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"token": token})
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials"})
}
