package general_api

import (
    
    "fmt"
    "errors"
    "net/http"
    "strings"
    jwt "github.com/dgrijalva/jwt-go"
    jwtmiddleware "github.com/auth0/go-jwt-middleware"

)



var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            secret := "empiricalfriedtreespracticalcarnies"
            if len(secret) == 0 {
                return nil, errors.New("Auth0 Client Secret Not Set")
            }
  
            return []byte(secret), nil
        },

    Extractor: extractTokenFromCookie,
    Debug: false,
})


func extractTokenFromCookie(r *http.Request) (string, error) {
    
    //get reference to cookie if set
    if cookie, err := r.Cookie("access_token"); err != nil {
    
        fmt.Println("Error reading cookie:", err)
        return "", err
        
    } else {
        
        auth := r.Header.Get("Authorization")

        if strings.Split(auth, " ")[1] != strings.Split(cookie.String()[13:], ".")[1] {  
            return "", errors.New("Token revoked, showed evidence of tampering.")
        }
        
        return cookie.String()[13:], nil
    }
}