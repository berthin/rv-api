package main

// @TODO: Fix identation and remove dead and testing-only code 
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
    "html"
	"net/http"
	"os"
	"time"
	//"io/ioutil"

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

// global
var db *sql.DB

func init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
                       DBUSER, DBPASSWORD, DBHOST, DBPORT, DBBASE)
    log.Println(dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if e := db.Ping(); e != nil {
		panic(err)
	}

	log.Println("Database connected!")
}

func main() {
	router := mux.NewRouter()

    router.HandleFunc("/", Index)

	router.HandleFunc("/api/v1/users", ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/{id:[0-9]+}", CreateUser).Methods(http.MethodGet)

	router.HandleFunc("/api/v1/widgets", ListWidgets).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/widgets", CreateWidget).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/widgets/{id:[0-9]+}", GetWidget).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/widgets/{id:[0-9]+}", UpdateWidget).Methods(http.MethodPut)

	router.HandleFunc("/api/v1/auth", GenToken).Methods(http.MethodGet)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Print("Listening on PORT " + PORT)
	log.Fatal(http.ListenAndServe(PORT, handlers.RecoveryHandler()(loggedRouter)))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// ****************
// Helper functions
// ****************

type Message struct {
	Message string `json:"message"`
}

func logErrorIfNonNil(err error) {
	if err != nil {
		log.Println(err)
	}
}

func sendMessage(message interface{}, statusCode int, w http.ResponseWriter) {
	msg, err := json.Marshal(message)
	logErrorIfNonNil(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(msg)
}

func verifyTokenAccess(w http.ResponseWriter, r *http.Request) bool {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
		
	if err != nil {
		sendMessage(Message{"Unauthorized access"}, http.StatusBadRequest, w)
		return false
	}
	
	if !token.Valid {
		sendMessage(Message{"Invalid token"}, http.StatusBadRequest, w)
		return false
	}
		
	return true
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

	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	authToken := Token{}
	authToken.Token = tokenString

	sendMessage(authToken, http.StatusOK, w)
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

	if !verifyTokenAccess(w, r) {
		return;
	}

	users := make([]User, 0)

	rows, err := db.Query("SELECT * FROM user ORDER BY user.id DESC")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {

		user := User{}

		err := rows.Scan(&user.Id, &user.Name, &user.Gravatar)
		logErrorIfNonNil(err)
		users = append(users, user)
	}

	logErrorIfNonNil(rows.Err())

	sendMessage(users, http.StatusOK, w)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if !verifyTokenAccess(w, r) {
		return;
	}

	user := User{}
	vars := mux.Vars(r)
	pathParam := vars["id"]

	log.Println("id" + pathParam)
	row := db.QueryRow("SELECT * FROM user WHERE user.id = ?", pathParam)

	err := row.Scan(&user.Id, &user.Name, &user.Gravatar)
	log.Println(user.Name)
	if err == sql.ErrNoRows {

		message := Message{"The record for the user couldn't be found."}
		sendMessage(message, http.StatusBadRequest, w)
		return
	}
	logErrorIfNonNil(err)

	sendMessage(user, http.StatusOK, w)
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
	if !verifyTokenAccess(w, r) {
		return
	}

	widgets := make([]Widget, 0)

	rows, err := db.Query("SELECT * FROM widget ORDER BY widget.id DESC")
	logErrorIfNonNil(err)
	defer rows.Close()

	for rows.Next() {

		widget := Widget{}

		err := rows.Scan(&widget.Id, &widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory)
		logErrorIfNonNil(err)

		widgets = append(widgets, widget)
	}

	logErrorIfNonNil(rows.Err())

	sendMessage(widgets, http.StatusOK, w)
}

// shortcut method
func GetWidget(w http.ResponseWriter, r *http.Request) {

	if !verifyTokenAccess(w, r) {
		return;
	}

	widget := Widget{}
	vars := mux.Vars(r)
	pathParam := vars["id"]

	row := db.QueryRow("SELECT * FROM widget WHERE widget.id = ?", pathParam)

	err := row.Scan(&widget.Id, &widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory)
	if err == sql.ErrNoRows {
		message := Message{"The record for the widget couldn't be found."}
		sendMessage(message, http.StatusBadRequest, w)
		return
	}

	logErrorIfNonNil(err)

	sendMessage(widget, http.StatusOK, w)
}

func CreateWidget(w http.ResponseWriter, r *http.Request) {
	if !verifyTokenAccess(w, r) {
		return;
	}

	widget := Widget{}

	err := json.NewDecoder(r.Body).Decode(&widget)
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
		//panic(err)
	//}
	//log.Println(string(body))
	//err = json.Unmarshal(body, &widget)

	if err != nil {
		log.Println(err)
		message := Message{"Something wrong happens. Try again later."}
		sendMessage(message, http.StatusBadRequest, w)
		return
	}
	
	if widget.Name == "" {
		log.Println(err) // DOES NOT MAKE SENSE TO HAVE THIS MESSAGE THOUGH
		message := Message{"Name is required."}
		sendMessage(message, http.StatusNotAcceptable, w)
		return
	}

	stmt, err := db.Prepare("INSERT INTO widget (name, color, price, melts, inventory) VALUES (?, ?, ?, ?, ?)")
	logErrorIfNonNil(err)
	defer stmt.Close()

	record, err := stmt.Exec(&widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory)
	logErrorIfNonNil(err)

	lastRecord, err := record.LastInsertId()
	logErrorIfNonNil(err)

	rowNum, err := record.RowsAffected()
	logErrorIfNonNil(err)

	log.Printf("Widget Created :: Record id %d :: Rows Affected %d", lastRecord, rowNum)

	message := Message{"Widget created successfully"}
	sendMessage(message, http.StatusOK, w)
}

func UpdateWidget(w http.ResponseWriter, r *http.Request) {

	if !verifyTokenAccess(w, r) {
		return;
	}
	
	widget := Widget{}
	vars := mux.Vars(r)
	pathParam := vars["id"]

	err := json.NewDecoder(r.Body).Decode(&widget)
	if err != nil {
		log.Println(err)
		message := Message{"Something wrong happens. Try again later."}
		sendMessage(message, http.StatusBadRequest, w)
		return
	}

	if widget.Name == "" {
		log.Println(err)
		message := Message{"Name is required."}
		sendMessage(message, http.StatusNotAcceptable, w)
		return
	}

	stmt, err := db.Prepare("UPDATE widget SET widget.name = ?, widget.color = ?, widget.price = ?, widget.melts = ?, widget.inventory = ? WHERE widget.id = ?")
	logErrorIfNonNil(err)
	defer stmt.Close()

	record, err := stmt.Exec(&widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory, pathParam)
	logErrorIfNonNil(err)

	rowNum, err := record.RowsAffected()
	logErrorIfNonNil(err)

	log.Printf("Widget Updated :: Rows Affected %d", rowNum)

	message := Message{"Widget updated successfully"}
	sendMessage(message, http.StatusOK, w)
}
