accord "test-accord" {
  source = "github.com/datascienceinc/accord//examples/test.hcl"
}


endpoint "/users" {
  method = "GET"

  response {
    headers {
      X-MY-HEADER = "stuff"
    }

    code = 300

    body = "hello world"
  }
}
