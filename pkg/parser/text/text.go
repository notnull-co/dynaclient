package text

import (
	"bytes"
	"io"

	"github.com/notnull-co/dynaclient/pkg/parser"
)

type textParser[T any] struct {
	textEncoder
}

type textEncoder struct {
}

func Encoder() parser.DynaEncoder {
	return &textEncoder{}
}

func New[T any]() parser.DynaParser[T] {
	return &textParser[T]{}
}

func (t *textParser[T]) Decode(reader io.ReadCloser) (*T, error) {
	return Decode[T](reader)
}

func (t *textEncoder) Encode(payload any) (io.ReadCloser, error) {
	return Encode(payload)
}

func Decode[T any](reader io.Reader) (*T, error) {
	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}
	var response any = string(bytes)
	assertedResponse := response.(T)

	return &assertedResponse, nil
}

func Encode(payload any) (io.ReadCloser, error) {
	if payload == nil {
		return nil, nil
	}
	payloadValue := *payload.(*any)

	return io.NopCloser(bytes.NewReader([]byte(payloadValue.(string)))), nil
}
