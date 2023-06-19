package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/notnull-co/dynaclient/pkg/parser"
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

	// payload := Payload{
	// 	Name: "muri",
	// 	Age:  10,
	// }

	req, err := NewRequest(http.MethodPost, svr.URL, "Oi")

	if err != nil {
		t.Fatal(err)
	}

	client := New[string]().WithParser(parser.Text)

	response, _, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}
}
