# gokit

A small collection of reusable Go utilities for building HTTP applications and other Go programs.

`gokit` provides lightweight helpers for environment variables, terminal logging, HTTP request logging, HTTP clients, JSON requests and responses, and input validation.

## Installation

```bash
go get github.com/modexusdev/gokit@latest
```

## Packages

### env

Simple environment variable handling with `.env` support.

```go
import "github.com/modexusdev/gokit/env"

port := env.Get("PORT", "8080")
```

If the environment variable does not exist or is empty, the fallback value is returned.

Example `.env`:

```env
PORT=9500
```

---

### crypto

Helpers for authenticated encryption and decryption using AES-256-GCM.

```go
import "github.com/modexusdev/gokit/crypto"
```

The package can be used for values that need to be stored encrypted but must later be recoverable, such as voucher codes, API credentials, or other application secrets.

#### Generate a key

`crypto` expects a Base64-encoded 32-byte key.

A suitable key can be generated with OpenSSL:

```bash
openssl rand -base64 32
```

Store the generated key securely, for example in an environment variable:

```env
ENCRYPTION_KEY=<generated-key>
```

Do not commit encryption keys to source control.

#### Encrypt a value

Use `Encrypt` to encrypt a string:

```go
key := env.Get(
    "ENCRYPTION_KEY",
    "",
)

encrypted, err := crypto.Encrypt(
    "SECRET-VALUE-123",
    key,
)
if err != nil {
    return err
}
```

The returned value is Base64 encoded and can safely be stored in a text field:

```text
YlJQY0...encrypted-data...
```

A random nonce is generated for every encryption operation. This means encrypting the same value multiple times produces different ciphertext.

#### Decrypt a value

Use `Decrypt` with the same key to recover the original value:

```go
decrypted, err := crypto.Decrypt(
    encrypted,
    key,
)
if err != nil {
    return err
}

fmt.Println(decrypted)
```

Output:

```text
SECRET-VALUE-123
```

The key used for decryption must be the same key that was used for encryption.

#### Empty values

Empty strings are returned unchanged:

```go
encrypted, err := crypto.Encrypt(
    "",
    key,
)
```

The result is:

```text
""
```

The same applies to `Decrypt`.

This is useful for optional encrypted fields.

#### Invalid keys and ciphertext

The encryption key must decode to exactly 32 bytes.

Invalid keys return:

```go
crypto.ErrInvalidKey
```

Invalid, corrupted, or modified encrypted values return:

```go
crypto.ErrInvalidCiphertext
```

Errors can be checked using `errors.Is`:

```go
value, err := crypto.Decrypt(
    encrypted,
    key,
)
if errors.Is(
    err,
    crypto.ErrInvalidCiphertext,
) {
    // Handle invalid or corrupted data.
}
```

AES-GCM provides authenticated encryption. If encrypted data is modified after encryption, decryption fails instead of returning corrupted plaintext.

#### Security considerations

Encryption keys should be treated as secrets.

They should:

- be generated using a cryptographically secure random source
- be stored outside source code
- not be committed to Git
- not be exposed to frontend applications
- not be written to application logs

Losing the encryption key means values encrypted with that key can no longer be recovered.

Changing the key also prevents existing encrypted values from being decrypted unless they are migrated to the new key first.

The package intentionally does not read environment variables itself. Applications provide the key explicitly, keeping `crypto` independent from configuration and allowing it to be used in HTTP servers, CLI applications, background workers, and other Go programs.

---

### log

Colored application and HTTP request logging.

```go
import "github.com/modexusdev/gokit/log"

logger := log.New("MyApp")

logger.Info("Server started")
logger.Warn("Something looks wrong")
logger.Error("Something failed")
logger.Debug("Debug information")
```

#### HTTP middleware

The logger can be used as standard `net/http` middleware.

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World"))
})

logger := log.New("MyApp")

handler := logger.HTTP(mux)

http.ListenAndServe(":8080", handler)
```

It logs information such as:

```text
[12:30:42] [MyApp] INFO GET / 200 15µs
```

HTTP status codes are automatically categorized:

| Status | Level |
|---|---|
| 2xx | INFO |
| 3xx | REDIRECT |
| 4xx | ERROR |
| 5xx | INTERNAL |

The middleware is based on `net/http` and can also be used with compatible routers such as Chi.

```go
r := chi.NewRouter()

r.Use(logger.HTTP)
```

#### HTTP log messages

An optional message can be attached to an HTTP request log using `log.Message`.

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

Messages are optional. If `log.Message` is not called, the request is logged normally without an additional message.

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

Log messages should not contain sensitive information such as passwords, authentication tokens, or other secrets.

---

### httpx

Helpers for common HTTP clients, requests, and responses.

```go
import "github.com/modexusdev/gokit/httpx"
```

#### HTTP client

`ClientSeconds` creates an `http.Client` with a timeout specified in seconds.

```go
client := httpx.ClientSeconds(5)
```

This is equivalent to:

```go
client := &http.Client{
    Timeout: 5 * time.Second,
}
```

The returned client is a standard `*http.Client` and can be used normally:

```go
res, err := client.Do(req)
if err != nil {
    return err
}
defer res.Body.Close()
```

#### Decode JSON requests

`DecodeJSON` decodes a JSON request body into a Go value.

```go
type createUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

var input createUserRequest

if err := httpx.DecodeJSON(r, &input); err != nil {
    httpx.Error(
        w,
        http.StatusBadRequest,
        "invalid request",
    )
    return
}
```

This reduces the common:

```go
json.NewDecoder(r.Body).Decode(&input)
```

pattern while still returning the decoding error to the caller.

#### Create JSON requests

`NewJSONRequest` creates an HTTP request with a JSON body.

```go
req, err := httpx.NewJSONRequest(
    ctx,
    http.MethodPost,
    "https://example.com/api/users",
    map[string]string{
        "name":  "John",
        "email": "john@example.com",
    },
)
if err != nil {
    return err
}
```

The helper:

- encodes the provided value as JSON
- creates the request using the provided context
- attaches the JSON body
- sets `Content-Type` to `application/json`

It replaces common boilerplate such as:

```go
data, err := json.Marshal(body)
if err != nil {
    return err
}

req, err := http.NewRequestWithContext(
    ctx,
    http.MethodPost,
    url,
    bytes.NewReader(data),
)
if err != nil {
    return err
}

req.Header.Set(
    "Content-Type",
    "application/json",
)
```

The returned value is a standard `*http.Request`, so additional headers can still be added normally:

```go
req.Header.Set(
    "Authorization",
    token,
)
```

#### OK

Returns a `200 OK` JSON response.

```go
httpx.OK(w, map[string]string{
    "status": "ok",
})
```

Response:

```json
{
    "status": "ok"
}
```

#### Error

Returns a JSON error response with a custom HTTP status code.

```go
httpx.Error(
    w,
    http.StatusBadRequest,
    "invalid request",
)
```

Response:

```json
{
    "error": "invalid request"
}
```

#### JSON

Use `JSON` when you need full control over the status code and response data.

```go
httpx.JSON(
    w,
    http.StatusServiceUnavailable,
    map[string]string{
        "status":  "unavailable",
        "service": "database",
    },
)
```

#### No Content

Returns `204 No Content`.

```go
httpx.NoContent(w)
```

## Example

```go
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

---

### validate

Lightweight validation helpers for strings, emails, integers, and common validation rules.

```go
import "github.com/modexusdev/gokit/validate"
```

Validators return `nil` when the value is valid and an error when a validation rule fails.

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

An empty description is allowed, while a non-empty description must not exceed 500 characters.

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

Validation is independent from `net/http`, but can easily be combined with `httpx`.

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

The validation package does not write HTTP responses itself. This keeps it reusable in HTTP servers, CLI applications, configuration loading, and other Go programs.

---

## Philosophy

`gokit` is intentionally small and built around Go's standard library.

The packages are designed to:

- reduce repetitive application code
- stay simple and explicit
- avoid unnecessary abstractions
- work with `net/http`
- remain framework-independent where possible

## License

MIT
