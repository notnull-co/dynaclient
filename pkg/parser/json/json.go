package json

import (
	"bytes"
	js "encoding/json"
	"io"

	"github.com/notnull-co/dynaclient/pkg/parser"
)

type jsonParser[T any] struct {
	jsonEncoder
}

type jsonEncoder struct {
}

func Encoder() parser.DynaEncoder {
	return &jsonEncoder{}
}

func New[T any]() parser.DynaParser[T] {
	return &jsonParser[T]{}
}

func (t *jsonParser[T]) Decode(reader io.ReadCloser) (*T, error) {
	return Decode[T](reader)
}

func (t *jsonEncoder) Encode(payload any) (io.ReadCloser, error) {
	return Encode(payload)
}

func Decode[T any](reader io.Reader) (*T, error) {
	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	response := new(T)

	if err := js.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func Encode(payload any) (io.ReadCloser, error) {
	b, err := js.Marshal(payload)

	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(b)), nil
}
