# gokit

A small collection of reusable Go utilities for building HTTP
applications.

`gokit` provides lightweight helpers for environment variables,
terminal logging, HTTP request logging, JSON responses, and input
validation.

## Installation

``` bash
go get github.com/modexusdev/gokit@latest
```

## Packages

### env

Simple environment variable handling with `.env` support.

``` go
import "github.com/modexusdev/gokit/env"

port := env.Get("PORT", "8080")
```

If the environment variable does not exist or is empty, the fallback
value is returned.

Example `.env`:

``` env
PORT=9500
```

------------------------------------------------------------------------

### log

Colored application and HTTP request logging.

``` go
import "github.com/modexusdev/gokit/log"

logger := log.New("MyApp")

logger.Info("Server started")
logger.Warn("Something looks wrong")
logger.Error("Something failed")
logger.Debug("Debug information")
```

#### HTTP middleware

The logger can be used as standard `net/http` middleware.

``` go
mux := http.NewServeMux()

mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World"))
})

logger := log.New("MyApp")

handler := logger.HTTP(mux)

http.ListenAndServe(":8080", handler)
```

It logs information such as:

``` text
[12:30:42] [MyApp] INFO GET / 200 15µs
```

HTTP status codes are automatically categorized:

  Status   Level
  -------- ----------
  2xx      INFO
  3xx      REDIRECT
  4xx      ERROR
  5xx      INTERNAL

The middleware is based on `net/http` and can also be used with
compatible routers such as Chi.

``` go
r := chi.NewRouter()

r.Use(logger.HTTP)
```

#### HTTP log messages

An optional message can be attached to an HTTP request log using
`log.Message`.

```go
r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
    log.Message(w, "Login successful")

    httpx.OK(w, map[string]string{
        "status": "ok",
    })
})
```

The message is automatically appended to the HTTP log:

```text
[12:30:42] [MyApp] INFO POST /login 200 15ms Login successful
```

Messages are optional. If `log.Message` is not called, the request is
logged normally without an additional message.

This can also be used for error responses:

```go
log.Message(w, "Invalid credentials")

httpx.Error(
    w,
    http.StatusUnauthorized,
    "invalid email or password",
)
```

Result:

```text
[12:30:42] [MyApp] ERROR POST /login 401 15ms Invalid credentials
```

Log messages should not contain sensitive information such as passwords,
authentication tokens, or other secrets.

------------------------------------------------------------------------

### httpx

Helpers for common HTTP responses.

``` go
import "github.com/modexusdev/gokit/httpx"
```

#### OK

Returns a `200 OK` JSON response.

``` go
httpx.OK(w, map[string]string{
    "status": "ok",
})
```

Response:

``` json
{
    "status": "ok"
}
```

#### Error

Returns a JSON error response with a custom HTTP status code.

``` go
httpx.Error(
    w,
    http.StatusBadRequest,
    "invalid request",
)
```

Response:

``` json
{
    "error": "invalid request"
}
```

#### JSON

Use `JSON` when you need full control over the status code and response
data.

``` go
httpx.JSON(
    w,
    http.StatusServiceUnavailable,
    map[string]string{
        "status": "unavailable",
        "service": "database",
    },
)
```

#### No Content

Returns `204 No Content`.

``` go
httpx.NoContent(w)
```

## Example

``` go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "github.com/modexusdev/gokit/env"
    "github.com/modexusdev/gokit/httpx"
    "github.com/modexusdev/gokit/log"
)

func main() {
    port := env.Get("PORT", "8080")

    logger := log.New("Example")

    r := chi.NewRouter()

    r.Use(logger.HTTP)

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        httpx.OK(w, map[string]string{
            "status": "ok",
        })
    })

    logger.Info("Server started on port", port)

    if err := http.ListenAndServe(":"+port, r); err != nil {
        logger.Error(err.Error())
    }
}
```

------------------------------------------------------------------------

### validate

Lightweight validation helpers for strings, emails, integers, and common
validation rules.

```go
import "github.com/modexusdev/gokit/validate"
```

Validators return `nil` when the value is valid and an error when a
validation rule fails.

#### Strings

Validate strings using optional rules.

```go
err := validate.String(
    "name",
    name,
    validate.Required(),
    validate.MinLength(2),
    validate.MaxLength(25),
)

if err != nil {
    fmt.Println(err)
}
```

Available string rules include:

- `Required()`
- `MinLength(n)`
- `MaxLength(n)`
- `OneOf(values...)`

For example, `OneOf` can restrict a value to a predefined set:

```go
err := validate.String(
    "role",
    role,
    validate.Required(),
    validate.OneOf(
        "user",
        "admin",
    ),
)
```

#### Email

Validate an email address with optional string rules.

```go
err := validate.Email(
    "email",
    email,
    validate.Required(),
    validate.MaxLength(254),
)
```

An invalid email returns a validation error.

```text
email must be a valid email address
```

#### Integers

Integer values can be validated using `Min` and `Max`.

```go
err := validate.Int(
    "port",
    port,
    validate.Min(1),
    validate.Max(65535),
)
```

Available integer rules include:

- `Min(n)`
- `Max(n)`

#### Optional values

Fields are optional unless `Required()` is specified.

For example:

```go
err := validate.String(
    "description",
    description,
    validate.MaxLength(500),
)
```

An empty description is allowed, while a non-empty description must not
exceed 500 characters.

To make the field required:

```go
err := validate.String(
    "description",
    description,
    validate.Required(),
    validate.MaxLength(500),
)
```

#### Using validation in an HTTP handler

Validation is independent from `net/http`, but can easily be combined
with `httpx`.

```go
if err := validate.Email(
    "email",
    input.Email,
    validate.Required(),
    validate.MaxLength(254),
); err != nil {
    httpx.Error(
        w,
        http.StatusBadRequest,
        err.Error(),
    )
    return
}
```

The validation package does not write HTTP responses itself. This keeps
it reusable in HTTP servers, CLI applications, configuration loading,
and other Go programs.

------------------------------------------------------------------------


## Philosophy

`gokit` is intentionally small and built around Go's standard library.

The packages are designed to:

-   reduce repetitive application code
-   stay simple and explicit
-   avoid unnecessary abstractions
-   work with `net/http`
-   remain framework-independent where possible

## License

MIT
