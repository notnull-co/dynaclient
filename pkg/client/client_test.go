package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/notnull-co/dynaclient/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestRequestString(t *testing.T) {
	type Response struct {
		Message string
	}

	response := Response{
		"You successfuly made it",
	}

	responseBytes, _ := json.Marshal(response)

	responseString := string(responseBytes)

	payloadChannel, svr := testServer(responseBytes, http.StatusOK)

	defer svr.Close()

	payload := "This is test data"

	req, err := NewRequest(http.MethodPost, svr.URL, payload)

	if assert.Nil(t, err) {
		client := New[string]().WithParser(parser.Text)

		response, httpResponse, err := client.Do(req)

		if assert.Nil(t, err) && assert.NotNil(t, response) {
			receivedPayload := <-payloadChannel
			assert.Equal(t, payload, string(receivedPayload))
			assert.Equal(t, responseString, *response)
			assert.Equal(t, responseString, httpResponse.BodyString())
		}
	}
}

func TestRequestJson(t *testing.T) {
	type Response struct {
		Message string
	}

	response := Response{
		"You successfuly made it",
	}

	responseBytes, _ := json.Marshal(response)

	payloadChannel, svr := testServer(responseBytes, http.StatusOK)

	defer svr.Close()

	type Payload struct {
		Name    string
		Age     int
		Balance float64
	}

	payload := Payload{
		Name:    "Muril Lo",
		Age:     69,
		Balance: 666.666,
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := NewRequest(http.MethodPost, svr.URL, payload)

	if assert.Nil(t, err) {
		client := New[Response]()

		response, httpResponse, err := client.Do(req)

		if assert.Nil(t, err) && assert.NotNil(t, response) {
			receivedPayload := <-payloadChannel

			var receveidParsedPayload Payload
			json.Unmarshal(receivedPayload, &receveidParsedPayload)

			assert.Equal(t, payload, receveidParsedPayload)
			assert.Equal(t, *response, *response)
			assert.Equal(t, responseBytes, httpResponse.Body())
			assert.Equal(t, string(payloadBytes), string(receivedPayload))
		}
	}
}

func TestRequestJsonWithFailConditionFailed(t *testing.T) {
	type Response struct {
		Message string
	}

	response := Response{
		"You successfuly made it",
	}

	responseBytes, _ := json.Marshal(response)

	_, svr := testServer(responseBytes, http.StatusBadRequest)

	defer svr.Close()

	type Payload struct {
		Name    string
		Age     int
		Balance float64
	}

	payload := Payload{
		Name:    "Muril Lo",
		Age:     69,
		Balance: 666.666,
	}

	req, err := NewRequest(http.MethodPost, svr.URL, payload)

	errBadRequest := errors.New("bad request")

	if assert.Nil(t, err) {
		client := New[Response]().WithFailCondition(func(r *http.Response) error {
			if r.StatusCode == http.StatusBadRequest {
				return errBadRequest
			}
			return nil
		})

		response, httpResponse, err := client.Do(req)

		if assert.NotNil(t, err) && assert.Error(t, errBadRequest, err) {
			assert.Nil(t, response)
			assert.NotNil(t, httpResponse)
			assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
		}
	}
}

func testServer(response []byte, statusCode int) (chan []byte, *httptest.Server) {
	payloadChannel := make(chan []byte, 1)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bytes, _ := io.ReadAll(r.Body)
		payloadChannel <- bytes
		w.WriteHeader(statusCode)
		w.Write(response)
	}))
	return payloadChannel, svr
}
