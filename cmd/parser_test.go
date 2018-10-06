package cmd

import (
	"net/http"
	"testing"
)

func TestParserParse(t *testing.T) {
	tests := []struct {
		ContentType string
		Body        interface{}
	}{
		{
			"application/json",
			`
				{
					"test": "value"
				}
			`,
		},
		{
			"application/xml",
			`
				<test>value</test>
			`,
		},
	}
	for _, test := range tests {
		headers := make(http.Header)
		headers.Add("Content-Type", test.ContentType)
		parser := Parser{Headers: headers, Body: test.Body}
		_, err := parser.Parse()
		if err != nil {
			t.Fatalf("Unexpected error: %e", err)
		}
	}
}
