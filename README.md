# Accord

Accord is a super simple consumer contract testing and stubbing tool, that is 
language agnostic.

## Getting Started

Prebuilt binaries can be downloaded [here](https://github.com/datascienceinc/accord/releases/latest),
however if you wish to build from source you will need `go 1.7+`

### Installing

```
go get github.com/datascienceinc/accord
```

### Writing a contract

A contract is defined by creating a `accord.hcl` file containing an "endpoint"
like the following.

```
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
```

this will define a stub/test-client at the url `/users` that will respond/request
with a `POST` method and the specified response/request data.

For the case of testing a provider with the "consumer" contract you can import
contracts from git, s3, http, local, and mercurial by adding the following to an
`accord.hcl`.

```
accord "test-accord" {
  source = "github.com/datascienceinc/accord//examples/test.hcl"
}
```
