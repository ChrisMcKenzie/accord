package accord

import (
	"github.com/hashicorp/hcl"
	"testing"
)

func TestParseHandlesQueryParams(t *testing.T) {
	// an accord config with a query entry in the request block
	config := `
    endpoint "/users" {
        method = "POST"

        request {

            query {
                hello   = "world"
                goodbye = "moon"
            }
        }
    }
    `
	// parse the string
	ast, err := hcl.Parse(config)
	// if there was something wrong
	if err != nil {
		t.Error(err.Error())
		panic(err)
	}

	// create the accord config from the ast
	contract, err := parse(ast)
	// if something went wrong
	if err != nil {
		t.Error(err.Error())
		panic(err)
	}
	// grab the query log from the config request
	query := contract.Endpoints[0].Request.Query
	// check the values of the query map
	if query["hello"] != "world" {
		// the test failed
		t.Errorf(
			"Incorrect value for hello in request query. Found %v, wanted %v",
			query["hello"], "world")
	}
}
