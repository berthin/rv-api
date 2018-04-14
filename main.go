package main


import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
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


func sendHttpMessage(message interface{}, statusCode int, w http.ResponseWriter) {
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
		sendHttpMessage(Message{"Unauthorized access"}, http.StatusBadRequest, w)
		return false
	}
	
	if !token.Valid {
		sendHttpMessage(Message{"Invalid token"}, http.StatusBadRequest, w)
		return false
	}
		
	return true
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

// ***************
// Auth
// ***************

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

	sendHttpMessage(authToken, http.StatusOK, w)
}


// ***************
// Users
// ***************

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
	sendHttpMessage(users, http.StatusOK, w)
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
		sendHttpMessage(message, http.StatusBadRequest, w)
		return
	}
	logOnError(err)

	sendHttpMessage(user, http.StatusOK, w)
}


// ***************
// Widgets
// ***************

type Widget struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Color     string  `json:"color,omitempty"`
	Price     float32 `json:"price"`
	Melts     bool    `json:"melts"`
	Inventory int64   `json:"inventory"`
}


func ListWidgets(w http.ResponseWriter, r *http.Request) {
	if !isTokenAccessCorrect(w, r) {
		return
	}

	rows, err := db.Query(`
		SELECT id, name, color, price, melts, inventory 
		FROM widget 
		ORDER BY widget.id DESC`)
	logOnError(err)
	defer rows.Close()

	widgets := make([]Widget, 0)
	for rows.Next() {
		widget := Widget{}

		err := rows.Scan(
			&widget.Id, 
			&widget.Name, 
			&widget.Color, 
			&widget.Price, 
			&widget.Melts, 
			&widget.Inventory)
		logOnError(err)

		widgets = append(widgets, widget)
	}
	logOnError(rows.Err())

	sendHttpMessage(widgets, http.StatusOK, w)
}


func GetWidget(w http.ResponseWriter, r *http.Request) {
	if !isTokenAccessCorrect(w, r) {
		return
	}

	vars := mux.Vars(r)
	pathParam := vars["id"]

	// Improve select to return in the column order fashion
	row := db.QueryRow(`
		SELECT id, name, color, price, melts, inventory 
		FROM widget 
		WHERE widget.id = ?`, pathParam) 

	widget := Widget{}
	err := row.Scan(
		&widget.Id,
		&widget.Name, 
		&widget.Color, 
		&widget.Price, 
		&widget.Melts, 
		&widget.Inventory)

	if err == sql.ErrNoRows {
		message := Message{"The record for the widget couldn't be found."}
		sendHttpMessage(message, http.StatusBadRequest, w)
		return
	}
	logOnError(err)

	sendHttpMessage(widget, http.StatusOK, w)
}


func CreateWidget(w http.ResponseWriter, r *http.Request) {
	if !isTokenAccessCorrect(w, r) {
		return;
	}

	widget := Widget{}

	err := json.NewDecoder(r.Body).Decode(&widget)

	if err != nil {
		logOnError(err)
		message := Message{"Something wrong happens. Try again later."}
		sendHttpMessage(message, http.StatusBadRequest, w)
		return
	}
	
	if widget.Name == "" {
		message := Message{"Name is required."}
		sendHttpMessage(message, http.StatusNotAcceptable, w)
		return
	}

	stmt, err := db.Prepare(`
		INSERT INTO widget 
		(name, color, price, melts, inventory) 
		VALUES (?, ?, ?, ?, ?)`)
	logOnError(err)
	defer stmt.Close()

	record, err := stmt.Exec(
		&widget.Name, 
		&widget.Color, 
		&widget.Price, 
		&widget.Melts, 
		&widget.Inventory)
	logOnError(err)

	lastRecord, err := record.LastInsertId()
	logOnError(err)

	rowNum, err := record.RowsAffected()
	logOnError(err)

	log.Printf("Widget Created :: Record id %d :: Rows Affected %d", lastRecord, rowNum)

	message := Message{"Widget created successfully"}
	sendHttpMessage(message, http.StatusOK, w)
}


func UpdateWidget(w http.ResponseWriter, r *http.Request) {
	if !isTokenAccessCorrect(w, r) {
		return
	}

	widget := Widget{}
	err := json.NewDecoder(r.Body).Decode(&widget)

	if err != nil {
		logOnError(err)
		message := Message{"Something wrong happens. Try again later."}
		sendHttpMessage(message, http.StatusBadRequest, w)
		return
	}

	if widget.Name == "" {
		message := Message{"Name is required."}
		sendHttpMessage(message, http.StatusNotAcceptable, w)
		return
	}

	stmt, err := db.Prepare(`
		UPDATE widget 
		SET widget.name = ?, 
			widget.color = ?,
			widget.price = ?, 
			widget.melts = ?, 
			widget.inventory = ?
		WHERE widget.id = ?`)
	logOnError(err)
	defer stmt.Close()

	vars := mux.Vars(r)
	pathParam := vars["id"]

	record, err := stmt.Exec(
		&widget.Name, 
		&widget.Color, 
		&widget.Price, 
		&widget.Melts, 
		&widget.Inventory, 
		pathParam)
	logOnError(err)

	rowNum, err := record.RowsAffected()
	logOnError(err)

	log.Printf("Widget Updated :: Rows Affected %d", rowNum)

	message := Message{"Widget updated successfully"}
	sendHttpMessage(message, http.StatusOK, w)
}