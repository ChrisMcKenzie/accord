
endpoint "/users" {
  method = "POST"

  request {
    body = <<-EOF
    {
      "test": "value"
    }
    EOF
  }

  response {
    code = 201

    body = "hello, world"
  }
}
