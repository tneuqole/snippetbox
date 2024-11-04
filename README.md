# snippetbox

This repo contains my code from reading the book [Let's Go](https://lets-go.alexedwards.net/)
by Alex Edwards.

## Running locally

Create the database user and database

```sql
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'password';

--TODO add database schema
```

Generate a TLS certificate

```zsh
mkdir tls
cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

Start the server from the root directory

```zsh
go run ./cmd/web
```
