// ***************
// Auth
// ***************

package main


import (
	"net/http"
	"time"

    jwt "github.com/dgrijalva/jwt-go"
    "github.com/dgrijalva/jwt-go/request"
    _ "github.com/go-sql-driver/mysql"
)


type Token struct {
    Token string `json:"token"`
}


var mySigningKey = []byte("4e3Fh54w374w")


func createToken(key []byte) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["admin"] = true
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    /* Sign the token with our secret */
    tokenString, err := token.SignedString(key)
    return tokenString, err
}


func GenToken(w http.ResponseWriter, r *http.Request) {
    tokenString, err := createToken(mySigningKey) 
    logOnError(err)
    authToken := Token{}
    authToken.Token = tokenString

    sendJsonThroughHttpMessage(authToken, http.StatusOK, w)
}


func isTokenAccessCorrect(w http.ResponseWriter, r *http.Request) bool {
    token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
        func(token *jwt.Token) (interface{}, error) {
            return mySigningKey, nil
        })
        
    if err != nil {
        message := Message{"Unauthorized access"}
        sendJsonThroughHttpMessage(message, http.StatusBadRequest, w)
        return false
    }
    
    if !token.Valid {
        message := Message{"Invalid token"}
        sendJsonThroughHttpMessage(message, http.StatusBadRequest, w)
        return false
    }
        
    return true
}