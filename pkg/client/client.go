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
	failCondition func(*http.Response) error
}

type DynaRequest struct {
	*http.Request
	parser.DynaEncoder
	payload any
}

type DynaResponse struct {
	*http.Response
	failed bool
}

func (d *DynaResponse) Body() []byte {
	bodyBytes, _ := io.ReadAll(d.Response.Body)
	return bodyBytes
}

func (d *DynaResponse) BodyString() string {
	bodyBytes := d.Body()

	if len(bodyBytes) > 0 {
		return string(bodyBytes)
	}

	return ""
}

func New[T any](client ...http.Client) *DynaClient[T] {
	var c http.Client

	if len(client) > 0 {
		c = client[0]
	}

	return &DynaClient[T]{
		c,
		json.New[T](),
		nil,
	}
}

func (c *DynaClient[T]) WithFailCondition(fn func(*http.Response) error) *DynaClient[T] {
	c.failCondition = fn
	return c
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

func (c *DynaClient[T]) Do(req *DynaRequest) (*T, *DynaResponse, error) {
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

	dynaResponse := DynaResponse{
		failed:   false,
		Response: response,
	}

	if err != nil {
		return nil, &dynaResponse, err
	}

	if c.failCondition != nil {
		if err := c.failCondition(dynaResponse.Response); err != nil {
			dynaResponse.failed = true
			return nil, &dynaResponse, err
		}
	}

	if response.Body != nil {
		bodyBytes, _ := io.ReadAll(dynaResponse.Response.Body)
		dynaResponse.Response.Body.Close()
		dynaResponse.Response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		responseValue, err := c.Decode(dynaResponse.Response.Body)

		if err != nil {
			return nil, &dynaResponse, err
		}

		return responseValue, &dynaResponse, nil
	}

	return nil, &dynaResponse, nil
}
