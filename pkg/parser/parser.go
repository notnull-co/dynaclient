package parser

import "io"

const (
	Json DynaParserType = 0
)

type DynaParserType int

type DynaParser[T any] interface {
	DynaEncoder
	DynaDecoder[T]
}

type DynaEncoder interface {
	Encode(any) (io.ReadCloser, error)
}

type DynaDecoder[T any] interface {
	Decode(io.ReadCloser) (*T, error)
}
