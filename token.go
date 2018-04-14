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


func GenToken(w http.ResponseWriter, r *http.Request) {
    token := jwt.New(jwt.SigningMethodHS256)

    // @TODO: Check what are these lines used for? Seems that claims is not used. By testing with
    // 		  GET, POST requests nothing changed when these were removed.
    claims := token.Claims.(jwt.MapClaims)
    claims["admin"] = true
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    /* Sign the token with our secret */
    tokenString, err := token.SignedString(mySigningKey)
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
        sendJsonThroughHttpMessage(Message{"Unauthorized access"}, http.StatusBadRequest, w)
        return false
    }
    
    if !token.Valid {
        sendJsonThroughHttpMessage(Message{"Invalid token"}, http.StatusBadRequest, w)
        return false
    }
        
    return true
}