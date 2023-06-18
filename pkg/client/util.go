package client

import (
	"github.com/notnull-co/dynaclient/pkg/parser"
	"github.com/notnull-co/dynaclient/pkg/parser/json"
)

var (
	encoderMap = map[parser.DynaParserType]parser.DynaEncoder{
		parser.Json: json.Encoder(),
	}
)

func getParser[T any](parserType parser.DynaParserType) parser.DynaParser[T] {
	switch parserType {
	case parser.Json:
		return json.New[T]()
	default:
		panic("invalid parser")
	}
}

func getEncoder(parserType parser.DynaParserType) parser.DynaEncoder {
	encoder, ok := encoderMap[parserType]

	if !ok {
		panic("invalid encoder")
	}

	return encoder
}
