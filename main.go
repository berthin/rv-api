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

    //http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        //fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    //})
    //log.Fatal(http.ListenAndServe(":8081", nil))

	log.Print("Listening on PORT " + PORT)
	log.Fatal(http.ListenAndServe(PORT, handlers.RecoveryHandler()(loggedRouter)))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// ***************
// Auth
// ***************

type Token struct {
	Token string `json:"token"`
}

type Message struct {
	Message string `json:"message"`
}

var mySigningKey = []byte("4e3Fh54w374w")

func GenToken(w http.ResponseWriter, r *http.Request) {

    log.Println("trying to get token")
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	authToken := Token{}
	authToken.Token = tokenString

	data, err := json.Marshal(authToken)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
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

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err == nil {
		if token.Valid {

			users := make([]User, 0)

			rows, err := db.Query("SELECT * FROM user ORDER BY user.id DESC")
			if err != nil {
				log.Println(err)
			}
			defer rows.Close()

			for rows.Next() {

				user := User{}

				err := rows.Scan(&user.Id, &user.Name, &user.Gravatar)
				if err != nil {
					log.Println(err)
				}
				users = append(users, user)
			}

			if err := rows.Err(); err != nil {
				log.Println(err)
			}

			usersData, err := json.Marshal(users)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(usersData)

		} else {
			message := Message{"Invalid token"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write(msg)
			return
		}
	} else {
		message := Message{"Unauthorized access"}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err == nil {
		if token.Valid {

			user := User{}
			vars := mux.Vars(r)
			pathParam := vars["id"]

            log.Println("id" + pathParam)
			row := db.QueryRow("SELECT * FROM user WHERE user.id = ?", pathParam)

			err := row.Scan(&user.Id, &user.Name, &user.Gravatar)
            log.Println(user.Name)
			if err == sql.ErrNoRows {

				message := Message{"The record for the user couldn't be found."}
				msg, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
				w.Write(msg)
				return
			}
			if err != nil {
				log.Println(err)
			}

			userData, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(userData)

		} else {
			message := Message{"Invalid token"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write(msg)
			return
		}
	} else {
		message := Message{"Unauthorized access"}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}
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

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err == nil {
		if token.Valid {

			widgets := make([]Widget, 0)

			rows, err := db.Query("SELECT * FROM widget ORDER BY widget.id DESC")
			if err != nil {
				log.Println(err)
			}
			defer rows.Close()

			for rows.Next() {

				widget := Widget{}

				err := rows.Scan(&widget.Id, &widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory)
				if err != nil {
					log.Println(err)
				}

				widgets = append(widgets, widget)
			}

			if err := rows.Err(); err != nil {
				log.Println(err)
			}

			widgetData, err := json.Marshal(widgets)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(widgetData)

		} else {
			message := Message{"Invalid token"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write(msg)
			return
		}
	} else {
		message := Message{"Unauthorized access"}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}
}

// shortcut method
func GetWidget(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err == nil {
		if token.Valid {

			widget := Widget{}
			vars := mux.Vars(r)
			pathParam := vars["id"]

			row := db.QueryRow("SELECT * FROM widget WHERE widget.id = ?", pathParam)

			err := row.Scan(&widget.Id, &widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory)
			if err == sql.ErrNoRows {

				message := Message{"The record for the widget couldn't be found."}
				msg, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
				w.Write(msg)
				return
			}
			if err != nil {
				log.Println(err)
			}

			widgetData, err := json.Marshal(widget)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(302)
			w.Write(widgetData)

		} else {
			message := Message{"Invalid token"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write(msg)
			return
		}
	} else {
		message := Message{"Unauthorized access"}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}
}

func CreateWidget(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err == nil {
		if token.Valid {

			widget := Widget{}

			err := json.NewDecoder(r.Body).Decode(&widget)
			if err != nil {
				log.Println(err)
				message := Message{"Something wrong happens. Try again later."}
				msg, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				w.Write(msg)
				return
			}
			if widget.Name == "" {
				log.Println(err)
				message := Message{"Name is required."}
				msg, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(406)
				w.Write(msg)
				return
			}

			stmt, err := db.Prepare("INSERT INTO widget (name, color, price, melts, inventory) VALUES (?, ?, ?, ?, ?)")
			if err != nil {
				log.Println(err)
			}
			defer stmt.Close()

			record, err := stmt.Exec(&widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory)
			if err != nil {
				log.Println(err)
			}

			lastRecord, err := record.LastInsertId()
			if err != nil {
				log.Println(err)
			}

			rowNum, err := record.RowsAffected()
			if err != nil {
				log.Println(err)
			}

			log.Printf("Widget Created :: Record id %d :: Rows Affected %d", lastRecord, rowNum)

			message := Message{"Widget created successfully"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(msg)

		} else {
			message := Message{"Invalid token"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write(msg)
			return
		}
	} else {
		message := Message{"Unauthorized access"}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}
}

func UpdateWidget(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err == nil {
		if token.Valid {

			widget := Widget{}
			vars := mux.Vars(r)
			pathParam := vars["id"]

			err := json.NewDecoder(r.Body).Decode(&widget)
			if err != nil {
				log.Println(err)
				message := Message{"Something wrong happens. Try again later."}
				msg, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				w.Write(msg)
				return
			}
			if widget.Name == "" {
				log.Println(err)
				message := Message{"Name is required."}
				msg, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(406)
				w.Write(msg)
				return
			}

			stmt, err := db.Prepare("UPDATE widget SET widget.name = ?, widget.color = ?, widget.price = ?, widget.melts = ?, widget.inventory = ? WHERE widget.id = ?")
			if err != nil {
				log.Println(err)
			}
			defer stmt.Close()

			record, err := stmt.Exec(&widget.Name, &widget.Color, &widget.Price, &widget.Melts, &widget.Inventory, pathParam)
			if err != nil {
				log.Println(err)
			}

			rowNum, err := record.RowsAffected()
			if err != nil {
				log.Println(err)
			}

			log.Printf("Widget Updated :: Rows Affected %d", rowNum)

			message := Message{"Widget updated successfully"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			w.Write(msg)

		} else {
			message := Message{"Invalid token"}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write(msg)
			return
		}
	} else {
		message := Message{"Unauthorized access"}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}
}
