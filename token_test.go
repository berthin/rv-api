// *******************
// Unit tests - token
// *******************

package main


import (
    "testing"
    "net/http"
    "net/http/httptest"
    "encoding/json"
    "fmt"
)


var router = configureRouter()


func getToken() (Token, error) {
    request, err := http.NewRequest("GET", "/api/v1/auth", nil)
    if err != nil {
        return Token{}, fmt.Errorf("Error while creating GET request. %s", err.Error())
    }

    response := httptest.NewRecorder()
    router.ServeHTTP(response, request)

    if response.Code != http.StatusOK {
        return Token{}, fmt.Errorf("Expected code %d, received %d.", http.StatusOK, response.Code)
    }

    received_token := Token{}
    err = json.NewDecoder(response.Body).Decode(&received_token)

    return received_token, err
}


func TestTokenAuth(t *testing.T) {
    token, err := getToken()
    validToken := token.Token

    if err != nil {
        t.Errorf("Error while decoding the token. %s.", err.Error())
    }
   
    // To test the token, we are sending a GET request with the received token 
    // and expect an OK as status code.
    request, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
        t.Errorf("Error while creating GET request. %s.", err.Error())
    }

    fakeToken1, err := createToken([]byte("abcdefg"))
    fakeToken2, err := createToken([]byte("invalid-token"))
    tokens := [] struct {
        Token string
        Code int
    }{
        {validToken, http.StatusOK},
        {fakeToken1, http.StatusBadRequest},
        {fakeToken2, http.StatusBadRequest},
    }

    for i := 0; i < 3; i++ {
        response := httptest.NewRecorder()
        request.Header.Set("authorization", tokens[i].Token)
        router.ServeHTTP(response, request)

        if response.Code != tokens[i].Code {
            t.Errorf("Invalid token! Expected status code %d, but received %d.", tokens[i].Code, response.Code)
        }
    }
}