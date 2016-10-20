package cmd

import (
	"fmt"
	"net/http"
	"time"

	accord "github.com/datascienceinc/accord/pkg"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the plan command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve accord as stub",
	Long: `
		generates a stub server based on the accord given for use by 
		consumers for development.
	`,
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue("Loaded the following endpoints from %s", cfgFile)
		color.Red("ERR: %s\n", serve())
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&port, "port", "p", "7600", "Port for server to listen on")
}

func serve() error {
	router := mux.NewRouter()

	ctx.ProcessEndpoints(func(ep *accord.Endpoint) {
		fmt.Printf("\n- ENDPOINT: [%s] %s\n", color.YellowString(ep.Method), color.BlueString(ep.URI))
		router.HandleFunc(ep.URI, newHandler(ep)).Methods(ep.Method)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()

}

func newHandler(ep *accord.Endpoint) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(ep.Response.Code)

		resp := parseBody(ep.Response.Headers, ep.Response.Body)
		w.Write(resp.Bytes())
	})
}
