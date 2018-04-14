package main


import (
    "encoding/json"
    "log"
    "net/http"

    jwt "github.com/dgrijalva/jwt-go"
    "github.com/dgrijalva/jwt-go/request"
    _ "github.com/go-sql-driver/mysql"
)


func panicOnError(err error) {
    if err != nil {
        panic(err)
    }
}


func logOnError(err error) {
    if err != nil {
        log.Println(err)
    }
}


func sendJsonThroughHttpMessage(message interface{}, statusCode int, w http.ResponseWriter) {
    msg, err := json.Marshal(message)
    logOnError(err)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(msg)
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