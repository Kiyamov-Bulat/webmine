package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"webmine/models"

	"github.com/dgrijalva/jwt-go"
)

var JwtAuth = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const TOKEN_NAME = "X-Session-Token"
		reqPath := r.URL.Path
		reqMethod := r.Method
		if isAllowedReq(reqMethod, reqPath) {
			next.ServeHTTP(w, r)
			return
		}
		tokenHeader := r.Header.Get(TOKEN_NAME)
		token, err := getToken(tokenHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := models.GetUser(token.UserId)
		if user.IsValid() {
			ctx := context.WithValue(r.Context(), "user", token.UserId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		} else {
			http.Error(w, "User is invalid", http.StatusMethodNotAllowed)
			return
		}
	})
}

func getToken(tokenHeader string) (*models.Token, error) {
	if tokenHeader != "" {
		splitted := strings.Split(tokenHeader, " ") //Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) == 2 {
			tokenPart := splitted[1]
			tk := &models.Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("AUTH_TOKEN")), nil
			})

			if err == nil {
				if token.Valid {
					log.Printf("User %d", tk.UserId)
					return tk, nil
				} else {
					return nil, errors.New("Token is not valid.")
				}
			} else {
				return nil, errors.New("Malformed authentication token")
			}
		} else {
			return nil, errors.New("Invalid/Malformed auth token")
		}
	}
	return nil, errors.New("Missing auth token")
}

func isAllowedReq(reqMethod string, reqPath string) bool {
	static := []string{"/static", "/data", "/frontend"}
	allowedPaths := []string{"/"}
	if reqMethod == "GET" {
		for _, val := range append(static, allowedPaths...) {
			if strings.HasPrefix(reqPath, val) {
				return true
			}
		}
	}
	if reqMethod == "POST" && reqPath == "/api/login" {
		return true
	}
	return false
}