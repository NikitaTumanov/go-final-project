package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("my_secret_key")

func createJWT() (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var signInRequest model.SignInRequest
	err := decoder.Decode(&signInRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errWrongPasswordJSON.Error(),
		})
		return
	}

	if os.Getenv("TODO_PASSWORD") != signInRequest.Password {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errWrongPassword.Error(),
		})
		return
	}

	jwt, err := createJWT()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errCreateJWT.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.SignInResponse{
		Token: jwt,
	})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("TODO_PASSWORD") != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "authentification required", http.StatusUnauthorized)
				return
			}

			cookieJWT := cookie.Value

			jwtToken, err := jwt.Parse(cookieJWT, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return secret, nil
			})

			if err != nil || !jwtToken.Valid {
				http.Error(w, "authentification failed", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
