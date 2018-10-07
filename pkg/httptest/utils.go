package httptest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	accord "github.com/ChrisMcKenzie/accord/pkg"
	"github.com/ChrisMcKenzie/accord/pkg/parser"
	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

func diffError(typ, a, b string) error {
	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: fmt.Sprintf("Actual %s", typ),
		ToFile:   fmt.Sprintf("Expected %s", typ),
		Context:  3,
		Eol:      "\n",
	}
	result, _ := difflib.GetContextDiffString(diff)
	return fmt.Errorf(strings.Replace(result, "\t", " ", -1))
}

func compareResponse(resp *http.Response, expect *accord.Response) error {
	var body bytes.Buffer
	defer resp.Body.Close()
	_, err := io.Copy(&body, resp.Body)
	if err != nil {
		color.Red("ERR: %s\n", err)
		return err
	}

	parser := parser.Parser{Headers: expect.Headers, Body: expect.Body}
	respBody, _ := parser.Parse()
	if body.String() != respBody.String() {
		diff := difflib.ContextDiff{
			A:        difflib.SplitLines(body.String()),
			B:        difflib.SplitLines(respBody.String()),
			FromFile: "Actual",
			ToFile:   "Expectation",
			Context:  3,
			Eol:      "\n",
		}
		result, _ := difflib.GetContextDiffString(diff)
		return fmt.Errorf(strings.Replace(result, "\t", " ", -1))
	}

	return nil
}
