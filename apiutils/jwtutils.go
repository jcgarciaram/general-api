package apiutils

import (
    jwt "github.com/dgrijalva/jwt-go"
    
    "time"
    "os"
)

var (
    secret = os.Getenv("AUTH0_CLIENT_SECRET")
)

func GenerateJWT(username string) (string, error) {
    
    mySigningKey := []byte(secret)

    type MyCustomClaims struct {
        Username string `json:"username"`
        jwt.StandardClaims
    }

    // Create the Claims
    claims := MyCustomClaims{
        username,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
            Issuer:    "referralapp",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(mySigningKey)
    
}