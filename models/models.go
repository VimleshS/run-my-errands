package models

import "github.com/dgrijalva/jwt-go"

const (
	//IndexRequestJob queue name
	IndexRequestJob = "RunAErrandQueue"
)

//User details, password is just a mock
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//JwtToken ..
type JwtToken struct {
	Token string `json:"token"`
}

// Exception ..
type Exception struct {
	Message string `json:"message"`
}

// JwtClaimsInfo ..
type JwtClaimsInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// Groceries ..
type Groceries struct {
	ID      int32     `json:"id"`
	Email   string    `json:"email"`
	Total   float64   `json:"total"`
	Items   []Grocery `json:"groceries"`
	Message string    `json:"message"`
}

// Grocery ..
type Grocery struct {
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
}
