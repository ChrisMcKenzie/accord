package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	accord "github.com/datascienceinc/accord/pkg"
	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

var client *http.Client

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

func init() {
	RootCmd.AddCommand(testCmd)
	client = &http.Client{}
}

func server(host, uri string) string {
	return fmt.Sprintf("%s%s", host, uri)
}

func test(host string) {
	ctx.ProcessEndpoints(func(ep *accord.Endpoint) {
		var buf bytes.Buffer
		if ep.Request != nil {
			buf = parseBody(ep.Request.Headers, ep.Request.Body)
		}

		req, err := http.NewRequest(ep.Method, server(host, ep.URI), &buf)
		if err != nil {
			color.Red("ERR: %s\n", err)
			return
		}

		res, err := client.Do(req)
		if err != nil {
			color.Red("ERR: %s\n", err)
			return
		}

		result := color.GreenString("OK")
		err = compareResponse(res, ep.Response)
		if err != nil {
			result = color.RedString("\nFAIL: \n%s\n", err.Error())
		}

		fmt.Printf("\n- ENDPOINT: [%s] %s | %s %s\n", color.YellowString(ep.Method), color.BlueString(ep.URI), res.Status, result)
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

func parseBody(h http.Header, i interface{}) bytes.Buffer {
	var buf bytes.Buffer
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
