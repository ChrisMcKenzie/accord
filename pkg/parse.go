package accord

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// LoadConfig ...
func LoadConfig(root string) (*Config, error) {
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

func parse(f *ast.File) (*Config, error) {
	// Top-level item should be the object list
	list, ok := f.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file does not contain root node object")
	}

	acc := new(Config)

	if endpoints := list.Filter("endpoint"); len(endpoints.Items) > 0 {
		var err error
		acc.Endpoints, err = loadEndpoints(endpoints)
		if err != nil {
			return nil, err
		}
	}

	if modules := list.Filter("accord"); len(modules.Items) > 0 {
		var err error
		acc.Modules, err = loadModules(modules)
		if err != nil {
			return nil, err
		}
	}

	return acc, nil
}

func loadModules(list *ast.ObjectList) ([]*Module, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*Module

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("accord '%s': should be an object", k)
		}

		var source string
		if o := listVal.Filter("source"); len(o.Items) > 0 {
			err := hcl.DecodeObject(&source, o.Items[0].Val)
			if err != nil {
				return nil, fmt.Errorf(
					"Error parsing source for %s: %s",
					k,
					err)
			}
		}

		result = append(result, &Module{k, source})
	}

	return result, nil
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
			return nil, fmt.Errorf("endpoint '%s': should be an object", uri)
		}

		var response *Response
		if o := listVal.Filter("response"); len(o.Items) > 0 {
			var err error
			response, err = loadResponse(o)
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
			var err error
			request, err = loadRequest(o)
			if err != nil {
				return nil, fmt.Errorf(
					"Error parsing request for %s: %s",
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

func loadResponse(list *ast.ObjectList) (*Response, error) {
	var result *Response
	item := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("request: should be an object")
	}

	var code int
	if o := listVal.Filter("code"); len(o.Items) > 0 {
		err := hcl.DecodeObject(&code, o.Items[0].Val)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing Response.Code for request: %s",
				err)
		}
	}

	var body interface{}
	if o := listVal.Filter("body"); len(o.Items) > 0 {
		err := hcl.DecodeObject(&body, o.Items[0].Val)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing Response.Body for request: %s",
				err)
		}
	}

	var headers http.Header
	if o := listVal.Filter("headers"); len(o.Items) > 0 {
		var err error
		headers, err = loadHeaders(o)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing Response.Headers: %s",
				err)
		}
	}

	result = &Response{
		Code:    code,
		Body:    body,
		Headers: headers,
	}

	return result, nil
}

func loadRequest(list *ast.ObjectList) (*Request, error) {
	var result *Request

	item := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("request: should be an object")
	}

	var body interface{}
	if o := listVal.Filter("body"); len(o.Items) > 0 {
		err := hcl.DecodeObject(&body, o.Items[0].Val)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing Request.Body for request: %s",
				err)
		}
	}

	headers := http.Header{}
	if o := listVal.Filter("headers"); len(o.Items) > 0 {
		var err error
		headers, err = loadHeaders(o)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing Request.Headers: %s",
				err)
		}
	}

	result = &Request{
		Body:    body,
		Headers: headers,
	}

	return result, nil
}

func loadHeaders(list *ast.ObjectList) (http.Header, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	result := http.Header{}

	item := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("request: should be an object")
	}

	var rawHeaders map[string]string
	if o := listVal.Filter("headers"); len(o.Items) > 0 {
		err := hcl.DecodeObject(&rawHeaders, o.Items[0].Val)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing Response.Code for request: %s",
				err)
		}
	}

	for k, val := range rawHeaders {
		result.Set(k, val)
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

func getModule(m *Module) ([]*Endpoint, error) {
	source, subDir := getter.SourceDirSubdir(m.Source)

	source, err := getter.Detect(source, "./.accord", getter.Detectors)
	if err != nil {
		return nil, fmt.Errorf("module %s: %s", m.Name, err)
	}

	// Check if the detector introduced something new.
	source, subDir2 := getter.SourceDirSubdir(source)
	if subDir2 != "" {
		subDir = filepath.Join(subDir2, subDir)
	}

	fmt.Println(source, subDir)

	return nil, err
}
