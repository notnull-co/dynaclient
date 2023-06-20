# DynaClient

DynaClient is a Go library that provides a custom HTTP client with support for dynamic types and parsers. It aims to have a minimal API and be very similar to the original Golang `http.Client` and `http.NewRequest`.

The library abstracts away the complexities of creating and handling HTTP requests, providing a streamlined and intuitive interface for making requests and processing responses.

## Installation

```shell
go get github.com/notnull-co/dynaclient
```

## Usage

First, import the necessary packages:

```go
import (
	"github.com/notnull-co/dynaclient/pkg/parser"
	"github.com/notnull-co/dynaclient/pkg/client"
)
```

### Creating a DynaClient

To create a DynaClient, use the `New` function:

```go
type ResponseDTO struct {
    Id string
    Username string
    Age int
}

c := client.New[ResponseDTO]()
```

You can optionally pass a custom `http.Client` to the `New` function if you want to use a specific HTTP client configuration.

### Setting the Parser

By default, the DynaClient uses the JSON parser provided by the library. You can also set a custom parser or switch between different parsers using the `WithCustomParser` and `WithParser` methods:

```go
c := client.New[ResponseDTO]().WithCustomParser(myCustomParser)
```

or

```go
c := client.New[ResponseDTO]().WithParser(parser.Yaml)
```

### Making Requests

To make a request with the DynaClient, create a `DynaRequest` using the `NewRequest` function. The `NewRequest` function is designed to be very similar to the original `http.NewRequest` function in Go's standard library:

```go
myStruct := MyPayloadStruct{
    Id: "abc-def-ghi",
    Book: "The tale of Muri",
    Price: 69,
}

req, err := client.NewRequest(http.HttpMethodPost, "https://api.example.com", myStruct)
if err != nil {
    // handle error
}
```

You can specify the request method, URL, an optional payload (body) for the request and a custom decoder. The DynaRequest object includes the underlying `http.Request` object and the appropriate `DynaEncoder` for the chosen parser.

### Sending the Request

To send the request and receive the response, use the `Do` method of the DynaClient. The `Do` method is designed to be similar to the `Do` method of the original `http.Client` in Go's standard library:

```go
parsedResponse, httpResponse, err := c.Do(req)
if err != nil {
    // handle error
}

// Access the response data
fmt.Println(response)
```

The `Do` method returns the decoded response data, an `http.Response` object, and an error (if any).

## Creating Custom Parsers

You can create custom parsers by implementing the `DynaParser` interface. The provided code includes an example implementation for a JSON parser (`jsonParser`). Feel free to modify and extend it according to your specific needs.

## Contributing

Contributions to DynaClient are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

DynaClient is licensed under the [Apache License](https://opensource.org/license/apache-2-0/).
