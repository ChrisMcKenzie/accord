package accord

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Load ...
func Load(root string) (*Accord, error) {
	data, err := readFile(root)
	if err != nil {
		return nil, err
	}

	f, err := hcl.Parse(data)
	if err != nil {
		return nil, err
	}

	acc, err := parse(f)

	return acc, err
}

func parse(f *ast.File) (*Accord, error) {
	// Top-level item should be the object list
	list, ok := f.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file does not contain root node object")
	}

	acc := new(Accord)

	if endpoints := list.Filter("endpoint"); len(endpoints.Items) > 0 {
		var err error
		acc.Endpoints, err = loadEndpoints(endpoints)
		if err != nil {
			return nil, err
		}
	}

	return acc, nil
}

func loadEndpoints(list *ast.ObjectList) ([]*Endpoint, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*Endpoint

	for _, item := range list.Items {
		uri := item.Keys[0].Token.Value().(string)

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("module '%s': should be an object", uri)
		}

		var response *Response
		if o := listVal.Filter("response"); len(o.Items) > 0 {
			err := hcl.DecodeObject(&response, o.Items[0].Val)
			if err != nil {
				return nil, fmt.Errorf(
					"Error parsing response for %s: %s",
					uri,
					err)
			}
		}

		var method string
		if o := listVal.Filter("method"); len(o.Items) > 0 {
			err := hcl.DecodeObject(&method, o.Items[0].Val)
			if err != nil {
				return nil, fmt.Errorf(
					"Error parsing response for %s: %s",
					uri,
					err)
			}
		}

		var request *Request
		if o := listVal.Filter("request"); len(o.Items) > 0 {
			err := hcl.DecodeObject(&request, o.Items[0].Val)
			if err != nil {
				return nil, fmt.Errorf(
					"Error parsing response for %s: %s",
					uri,
					err)
			}
		}

		result = append(result, &Endpoint{
			URI:      uri,
			Method:   method,
			Request:  request,
			Response: response,
		})
	}

	return result, nil
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
