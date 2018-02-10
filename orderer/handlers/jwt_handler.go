package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/VimleshS/run-my-errands/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

//CreateTokenEndpoint Authenticates user with email and password and generate hmac token valid for 15 minutes
func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(models.Exception{Message: err.Error()})
		logrus.WithField("user", user).Info("Error decoding user")
		return
	}

	// Create the Claims
	claims := models.JwtClaimsInfo{
		user.Email,
		user.Password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "run-my-errand",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, error := token.SignedString([]byte(models.SignInSecret))
	if error != nil {
		json.NewEncoder(w).Encode(models.Exception{Message: err.Error()})
		logrus.WithField("token", token).Info("Error signing keys")
		return
	}
	json.NewEncoder(w).Encode(models.JwtToken{Token: tokenString})
}

//ValidateToken it acts as a middleware validating the token and pipeling request further
func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(models.SignInSecret), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(models.Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "user", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(models.Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(models.Exception{Message: "An authorization header is required"})
		}
	})
}
