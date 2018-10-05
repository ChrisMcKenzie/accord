# Accord

[![wercker status](https://app.wercker.com/status/713a6a1f73487050f80c6747150b9e4e/s/master "wercker status")](https://app.wercker.com/project/byKey/713a6a1f73487050f80c6747150b9e4e)

Accord is a super simple consumer contract testing and stubbing tool, that is 
language agnostic.

## Getting Started

Prebuilt binaries can be downloaded [here](https://github.com/ChrisMcKenzie/accord/releases/latest),
however if you wish to build from source you will need `go 1.7+`

### Installing

```
go get github.com/ChrisMcKenzie/accord
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
  source = "github.com/ChrisMcKenzie/accord//examples/test.hcl"
}
```

You can create multiple `accord` definitions and they will all be pulled and merged
in to a single suite for serving and testing.

## Running A Stub Server

Accord can be used to run a stub/mock server based on the endpoints in your 
`accord.hcl`. This can be done by running the following.

```
accord serve
Loaded the following endpoints from examples/accord.hcl

Module root:
        ENDPOINT: [GET] /users
Module test-accord:
        ENDPOINT: [POST] /users
```

this will start a web server listening on `localhost:7600` the port can be 
changed by using the `-p` flag.
