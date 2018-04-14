// ***************
// Users
// ***************

package main


import (
    "database/sql"
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
)


type User struct {
    Id       int64  `json:"id"`
    Name     string `json:"name"`
    Gravatar string `json:"gravatar,omitempty"`
}


func ListUsers(w http.ResponseWriter, r *http.Request) {
    if !isTokenAccessCorrect(w, r) {
        return
    }

    rows, err := db.Query(`
        SELECT id, name, gravatar
        FROM user 
        ORDER BY user.id DESC`)
    logOnError(err)
    defer rows.Close()

    users := make([]User, 0)
    for rows.Next() {
        user := User{}

        err := rows.Scan(
            &user.Id, 
            &user.Name, 
            &user.Gravatar)
        logOnError(err)

        users = append(users, user)
    }

    logOnError(rows.Err())
    sendJsonThroughHttpMessage(users, http.StatusOK, w)
}


func ListUserById(w http.ResponseWriter, r *http.Request) {
    if !isTokenAccessCorrect(w, r) {
        return
    }

    vars := mux.Vars(r)
    pathParam := vars["id"]

    log.Println("id" + pathParam)
    row := db.QueryRow(`
        SELECT id, name, gravatar 
        FROM user 
        WHERE user.id = ?`, pathParam)

    user := User{}
    err := row.Scan(
        &user.Id, 
        &user.Name, 
        &user.Gravatar)
    log.Println(user.Name)

    if err == sql.ErrNoRows {
        message := Message{"The record for the user couldn't be found."}
        sendJsonThroughHttpMessage(message, http.StatusBadRequest, w)
        return
    }
    logOnError(err)

    sendJsonThroughHttpMessage(user, http.StatusOK, w)
}