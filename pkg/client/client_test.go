package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Payload struct {
	Name string
	Age  int
}

type Response struct {
	Message string
}

func TestClient(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "Message": "You successfuly made it" }`)
	}))
	defer svr.Close()

	payload := Payload{
		Name: "muri",
		Age:  10,
	}

	req, err := NewRequest(http.MethodPost, svr.URL, payload)

	if err != nil {
		t.Fatal(err)
	}

	response, _, err := New[Response]().Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}
}
