package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/notnull-co/dynaclient/pkg/parser"
	"github.com/notnull-co/dynaclient/pkg/parser/json"
)

type DynaClient[T any] struct {
	http.Client
	parser.DynaParser[T]
}

type DynaRequest struct {
	*http.Request
	parser.DynaEncoder
	payload any
}

func New[T any](client ...http.Client) *DynaClient[T] {
	var c http.Client

	if len(client) > 0 {
		c = client[0]
	}

	return &DynaClient[T]{
		c,
		json.New[T](),
	}
}

func (c *DynaClient[T]) WithCustomParser(parser parser.DynaParser[T]) *DynaClient[T] {
	c.DynaParser = parser
	return c
}

func (c *DynaClient[T]) WithParser(parserType parser.DynaParserType) *DynaClient[T] {
	c.DynaParser = getParser[T](parserType)
	return c
}

func NewRequest(method string, url string, body interface{}, parserType ...parser.DynaParserType) (*DynaRequest, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	var encoder parser.DynaEncoder
	if len(parserType) > 0 {
		encoder = getEncoder(parserType[0])
	}

	return &DynaRequest{
		Request:     req,
		DynaEncoder: encoder,
		payload:     &body,
	}, nil
}

func (c *DynaClient[T]) Do(req *DynaRequest) (*T, *http.Response, error) {

	encoder := c.DynaParser.Encode
	if req.DynaEncoder != nil {
		encoder = req.DynaEncoder.Encode
	}

	body, err := encoder(req.payload)

	if err != nil {
		return nil, nil, err
	}

	req.Body = body

	response, err := c.Client.Do(req.Request)

	if err != nil {
		return nil, response, err
	}

	if response.Body != nil {
		bodyBytes, _ := io.ReadAll(response.Body)
		response.Body.Close()
		response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		responseValue, err := c.Decode(response.Body)

		if err != nil {
			return nil, response, err
		}

		return responseValue, response, nil
	}

	return nil, response, nil
}
