package controllers

import (
    "fmt"
    "net/http"
    "testing"
)

func TestGetBearerTokenValidHeader(t *testing.T) {
    expectedToken := "testing"

    r := &http.Request{
        Header: make(http.Header),
    }
    r.Header.Add("Authentication", fmt.Sprintf("Bearer %s", expectedToken))

    actualToken, err := GetBearerToken(r)
    if err != nil {
        t.Fatal("Error retrieving bearer token")
    }
    if actualToken != expectedToken {
        t.Errorf("Access token '%s' does not match '%s'", actualToken, expectedToken)
    }
}
