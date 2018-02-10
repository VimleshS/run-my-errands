package models

import "github.com/dgrijalva/jwt-go"

const (
	//IndexRequestJob queue name
	IndexRequestJob = "RunAErrandQueue"
	//SignInSecret Used for siging token
	SignInSecret = "secret"
)

//User It holds the users information
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//JwtToken It holds a serializable Token
type JwtToken struct {
	Token string `json:"token"`
}

//Exception It holds a serializable Exception
type Exception struct {
	Message string `json:"message"`
}

//JwtClaimsInfo is used by middleware, to pass information
type JwtClaimsInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//Groceries It holds a serializable Grocery Items
type Groceries struct {
	ID      int32     `json:"id"`
	Email   string    `json:"email"`
	Total   float64   `json:"total"`
	Items   []Grocery `json:"groceries"`
	Message string    `json:"message"`
}

//Grocery Holds grocery item detail
type Grocery struct {
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
}
