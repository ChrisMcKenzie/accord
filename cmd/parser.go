package cmd

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type BodyParser interface {
	Parse() (byteBufferReadCloser, error)
}

type Parser struct {
	Headers http.Header
	Body    interface{}
}

func (parser Parser) Parse() (byteBufferReadCloser, error) {
	var buf byteBufferReadCloser
	if parser.Body == nil {
		parser.Body = ""
	}
	switch parser.Headers.Get("Content-Type") {
	case "application/json":
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		err := enc.Encode(parser.Body)
		if err != nil {
			return byteBufferReadCloser{}, err
		}
	case "application/xml":
		enc := xml.NewEncoder(&buf)
		err := enc.Encode(parser.Body)
		if err != nil {
			return byteBufferReadCloser{}, err
		}
	default:
		buf.WriteString(parser.Body.(string))
	}
	return buf, nil
}
