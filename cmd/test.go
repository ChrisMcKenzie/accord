package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	accord "github.com/datascienceinc/accord/pkg"
	"github.com/datascienceinc/accord/pkg/httptest"
	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

var client *httptest.Client

// testCmd represents the plan command
var testCmd = &cobra.Command{
	Use:    "test URL",
	Short:  "test accord against server",
	Long:   ``,
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			color.Red("A server URL is required.")
			return
		}
		color.Blue("Loaded the following endpoints from %s\n", cfgFile)

		test(args[0])
	},
}

type byteBufferReadCloser struct {
	bytes.Buffer
}

func (b *byteBufferReadCloser) Close() error {
	return nil
}

func init() {
	RootCmd.AddCommand(testCmd)
	client = httptest.NewClient()
}

func server(host, uri string, query map[string]string) *url.URL {

	url := url.URL{Host: host, Path: uri}
	// if there are query parameters to add to the url
	if len(query) != 0 {
		// grab the query object
		q := url.Query()
		// loop over the endpoint query specifications
		for k, v := range query {
			// add the query to the url
			q.Add(k, v)
		}
		// assign the query back to the url
		url.RawQuery = q.Encode()
	}

	return &url
}

func test(host string) {
	ctx.ProcessEndpoints(func(ep *accord.Endpoint) {
		var buf byteBufferReadCloser
		if ep.Request != nil {
			buf = parseBody(ep.Request.Headers, ep.Request.Body)
		}

		req := &http.Request{
			URL:    server(host, ep.URI, ep.Request.Query),
			Method: ep.Method,
			Body:   &buf,
		}

		err := client.Evaluate(req, ep.Response)
		if err != nil {
			color.Red("ERR: %s\n", err)
			return
		}

		result := color.GreenString("OK")
		if err != nil {
			result = color.RedString("\nFAIL: \n%s\n", err.Error())
		}

		fmt.Printf("\n- ENDPOINT: [%s] %s | %s\n", color.YellowString(ep.Method), color.BlueString(ep.URI), result)
	})
}

func compareResponse(resp *http.Response, expect *accord.Response) error {
	var body bytes.Buffer
	defer resp.Body.Close()
	_, err := io.Copy(&body, resp.Body)
	if err != nil {
		color.Red("ERR: %s\n", err)
		return err
	}

	respBody := parseBody(expect.Headers, expect.Body)
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

func parseBody(h http.Header, i interface{}) byteBufferReadCloser {
	var buf byteBufferReadCloser
	if i == nil {
		i = ""
	}

	if _, ok := i.(string); h.Get("Content-Type") == "application/json" || !ok {
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		enc.Encode(i)
	} else {
		buf.WriteString(i.(string))
	}

	return buf
}
