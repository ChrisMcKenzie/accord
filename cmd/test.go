package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	accord "github.com/datascienceinc/accord/pkg"
	"github.com/fatih/color"
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
		color.Blue("Loaded the following endpoints from %s\n\n", cfgFile)

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
		var buf io.Reader
		if ep.Request != nil && ep.Request.Body != "" {
			buf = strings.NewReader(ep.Request.Body)
		}
		req, err := http.NewRequest(ep.Method, server(host, ep.URI), buf)
		if err != nil {
			color.Red("ERR: %s\n", err)
			return
		}

		for header, value := range ep.Response.Headers {
			req.Header.Add(header, value)
		}

		res, err := client.Do(req)
		if err != nil {
			color.Red("ERR: %s\n", err)
			return
		}

		var body bytes.Buffer
		defer res.Body.Close()
		_, err = io.Copy(&body, res.Body)
		if err != nil {
			color.Red("ERR: %s\n", err)
			return
		}

		result := color.GreenString("OK")
		if body.String() != ep.Response.Body {
			result = color.RedString("Fail %s", body.String())
		}

		fmt.Printf("\tENDPOINT: [%s] %s | %s %s\n", color.YellowString(ep.Method), color.BlueString(ep.URI), res.Status, result)
	})
}
