accord "test-accord" {
  source = "github.com/datascienceinc/accord//examples/test.hcl"
}


endpoint "/users" {
  method = "GET"

  response {
    headers {
      X-MY-HEADER = "stuff"
			Content-Type = "application/json"
    }

    code = 300

    body {
			message = "hello world"
		}
  }
}
