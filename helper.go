// ****************
// Helper functions
// ****************

package main


import (
    "encoding/json"
    "log"
    "net/http"
)


type Message struct {
    Message string `json:"message"`
}


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