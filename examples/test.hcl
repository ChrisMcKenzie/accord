
endpoint "/users" {
  method = "POST"

  request {
    body = <<-EOF
    {
      "test": "value"
    }
    EOF

    query {
      hello   = "world"
      goodbye = "moon"
    }
  }

  response {
    code = 201

    body = "hello, world"
  }
}
