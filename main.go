package main


import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)


const (
    DBUSER     = "root"
    DBPASSWORD = "1234"
    DBHOST     = "127.0.0.1"
    DBPORT     = "8080"
    DBBASE     = "redventures"
    PORT       = ":8081"
)


// ****************
// Helper functions
// ****************

type Message struct {
    Message string `json:"message"`
}


// global
var db *sql.DB


func init() {
    var err error
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
                       DBUSER, DBPASSWORD, DBHOST, DBPORT, DBBASE)
    log.Println(dsn)

    db, err = sql.Open("mysql", dsn)
    panicOnError(err)

    panicOnError(db.Ping())

    log.Println("Database connected!")
}


func main() {
    router := mux.NewRouter()

    router.HandleFunc("/api/v1/users",
                      ListUsers).Methods(http.MethodGet)
    router.HandleFunc("/api/v1/users/{id:[0-9]+}", 
                      ListUserById).Methods(http.MethodGet)

    router.HandleFunc("/api/v1/widgets", 
                      ListWidgets).Methods(http.MethodGet)
    router.HandleFunc("/api/v1/widgets", 
                      CreateWidget).Methods(http.MethodPost)
    router.HandleFunc("/api/v1/widgets/{id:[0-9]+}", 
                      GetWidget).Methods(http.MethodGet)
    router.HandleFunc("/api/v1/widgets/{id:[0-9]+}", 
                      UpdateWidget).Methods(http.MethodPut)

    router.HandleFunc("/api/v1/auth", 
                      GenToken).Methods(http.MethodGet)

    loggedRouter := handlers.LoggingHandler(os.Stdout, router)

    log.Print("Listening on PORT " + PORT)
    log.Fatal(http.ListenAndServe(PORT, handlers.RecoveryHandler()(loggedRouter)))
}