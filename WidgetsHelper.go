// ***************
// Widgets
// ***************

package main


import (
    "database/sql"
    "log"
	"net/http"
    "encoding/json"
	
    "github.com/gorilla/mux"
)


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

    sendJsonThroughHttpMessage(widgets, http.StatusOK, w)
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
        sendJsonThroughHttpMessage(message, http.StatusBadRequest, w)
        return
    }
    logOnError(err)

    sendJsonThroughHttpMessage(widget, http.StatusOK, w)
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
        sendJsonThroughHttpMessage(message, http.StatusBadRequest, w)
        return
    }
    
    if widget.Name == "" {
        message := Message{"Name is required."}
        sendJsonThroughHttpMessage(message, http.StatusNotAcceptable, w)
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
    sendJsonThroughHttpMessage(message, http.StatusOK, w)
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
        sendJsonThroughHttpMessage(message, http.StatusBadRequest, w)
        return
    }

    if widget.Name == "" {
        message := Message{"Name is required."}
        sendJsonThroughHttpMessage(message, http.StatusNotAcceptable, w)
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
    sendJsonThroughHttpMessage(message, http.StatusOK, w)
}