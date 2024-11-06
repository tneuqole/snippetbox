# snippetbox

This repo contains my code from reading the book [Let's Go](https://lets-go.alexedwards.net/)
by Alex Edwards.

## Running the server locally

See internal/models/testdata for database setup.

Generate a TLS certificate.

```zsh
mkdir tls
cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

Start the server from the root directory.

```zsh
go run ./cmd/web
```

## Running tests

To run all tests and view coverage:

```zsh
go test --covermode=count -coverprofile=profile.out ./...
go tool cover -html=profile.out
```
