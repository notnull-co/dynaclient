package xml

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/notnull-co/dynaclient/pkg/parser"
)

type xmlParser[T any] struct {
	xmlEncoder
}

type xmlEncoder struct {
}

func Encoder() parser.DynaEncoder {
	return &xmlEncoder{}
}

func New[T any]() parser.DynaParser[T] {
	return &xmlParser[T]{}
}

func (t *xmlParser[T]) Decode(reader io.ReadCloser) (*T, error) {
	return Decode[T](reader)
}

func (t *xmlEncoder) Encode(payload any) (io.ReadCloser, error) {
	return Encode(payload)
}

func Decode[T any](reader io.Reader) (*T, error) {
	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	response := new(T)

	if err := xml.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func Encode(payload any) (io.ReadCloser, error) {
	if payload == nil {
		return nil, nil
	}

	b, err := xml.Marshal(payload)

	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(b)), nil
}
