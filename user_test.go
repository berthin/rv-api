// *******************
// Unit tests - User
// *******************

package main


import (
    "testing"
    "net/http"
    "net/http/httptest"
    "encoding/json"
    "fmt"
)


const (
    // `MAX_GET_USER_REQUESTS` defines the number of GET requests
    // that are going to be tested.
    MAX_GET_USER_REQUESTS = 5
)


type ArrayOfUsers struct {
    Users []User `json:"array"`
}


func testGetUserByID(token string, userId int64) error {
    userIdAsString := fmt.Sprintf("%d", userId)
    request, err := http.NewRequest("GET", "/api/v1/users/" + userIdAsString, nil)
    if err != nil {
        fmt.Errorf("Error while creating GET request. %s.", err.Error())
    }

    response := httptest.NewRecorder()
    request.Header.Set("authorization", token)
    routerForUnitTests.ServeHTTP(response, request)

    if response.Code != http.StatusOK {
        return fmt.Errorf("Invalid token! Expected status code %d, but received %d.", 
                          http.StatusOK, response.Code)
    }
    return nil
}


func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}


func TestGetUsers(t *testing.T) {
    token, err := getToken()
    validToken := token.Token

    if err != nil {
        t.Errorf("Error while decoding the token. %s.", err.Error())
    }
   
    request, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
        t.Errorf("Error while creating GET request. %s.", err.Error())
    }

    response := httptest.NewRecorder()
    request.Header.Set("authorization", validToken)
    routerForUnitTests.ServeHTTP(response, request)

    if response.Code != http.StatusOK {
        t.Errorf("Invalid token! Expected status code %d, but received %d.", 
                 http.StatusOK, response.Code)
    }
    
    usersJson := ArrayOfUsers{}
    err = json.NewDecoder(response.Body).Decode(&usersJson.Users)
    users := usersJson.Users

    if err != nil {
        t.Errorf("Error when parsing the retorned list of `users`. %s", err.Error())
    }

    if numUsers := min(len(users), MAX_GET_USER_REQUESTS); numUsers > 0 {
        for i := 0; i < numUsers; i++ {
            err = testGetUserByID(validToken, users[i].Id)
            if err != nil {
                t.Errorf("Error when testing GET single user with id: %d. %s", 
                         users[i].Id, err.Error())
            }
        } 
    }
}