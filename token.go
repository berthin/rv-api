// ***************
// Auth
// ***************

package main


import (
	"net/http"
	"time"

    jwt "github.com/dgrijalva/jwt-go"
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