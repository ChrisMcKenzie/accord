package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type ByteBufferReadCloser struct {
	bytes.Buffer
}

func (b *ByteBufferReadCloser) Close() error {
	return nil
}

type BodyParser interface {
	Parse() (ByteBufferReadCloser, error)
}

type Parser struct {
	Headers http.Header
	Body    interface{}
}

func (parser Parser) Parse() (ByteBufferReadCloser, error) {
	var buf ByteBufferReadCloser
	if parser.Body == nil {
		parser.Body = ""
	}
	switch parser.Headers.Get("Content-Type") {
	case "application/json":
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		err := enc.Encode(parser.Body)
		if err != nil {
			return ByteBufferReadCloser{}, err
		}
	case "application/xml":
		enc := xml.NewEncoder(&buf)
		err := enc.Encode(parser.Body)
		if err != nil {
			return ByteBufferReadCloser{}, err
		}
	default:
		buf.WriteString(parser.Body.(string))
	}
	return buf, nil
}
