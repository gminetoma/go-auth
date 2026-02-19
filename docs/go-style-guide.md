# Go Style Guide

Based on the [Google Go Style Guide](https://google.github.io/styleguide/go/).

## Principles (in order of priority)

1. **Clarity** — code purpose and rationale should be evident to readers
2. **Simplicity** — accomplish goals using the simplest approach
3. **Concision** — maintain high signal-to-noise ratio
4. **Maintainability** — enable future modifications without difficulty
5. **Consistency** — align with existing codebase patterns

## Formatting

- All code must pass `gofmt`
- Aim for 80-column comments; no hard line length limit for code
- Prefer keeping function signatures on a single line
- Use `goimports` for import ordering

## Imports

Group imports in this order, separated by blank lines:

```go
import (
    // Standard library
    "context"
    "fmt"

    // Third-party
    "github.com/golang-jwt/jwt/v5"

    // Local packages
    "github.com/gminetoma/go-auth/src/core/credentials/domain"
    "github.com/gminetoma/go-shared/src/errs"
)
```

## Naming

### General

- Use `MixedCaps` or `mixedCaps`. Never use `snake_case` or `SCREAMING_SNAKE_CASE`
- Exported: `PascalCase`. Unexported: `camelCase`
- Name length should be proportional to scope size

```go
// Tight scope (few lines) — short names are clear
for i, v := range items {
    process(v, i)
}

rt, err := r.queries.GetRefreshTokenByToken(ctx, token)
if err != nil {
    return nil, err
}

return &domain.RefreshToken{ID: domain.RefreshTokenID(rt.ID)}, nil

// Wider scope (used across many lines) — descriptive names
credentialsRepo := credentialsInfra.NewPGCredentialsRepository(...)
userRepo := userInfra.NewPGUserRepository(...)
// ... 20 lines later ...
if err := credentialsRepo.Create(ctx, creds); err != nil { ... }

// Package-level / exported — always descriptive
var ErrInvalidCredentials = errors.New("auth.invalid-credentials")
type AuthService struct { ... }
```

### Import Aliases

Use full feature name + layer to avoid ambiguity across files:

```go
// Good
credentialsApplication "github.com/gminetoma/go-auth/src/core/credentials/application"
credentialsDomain "github.com/gminetoma/go-auth/src/core/credentials/domain"
userApplication "github.com/gminetoma/go-auth/src/core/user/application"

// Bad — abbreviated feature names
credsApplication "github.com/gminetoma/go-auth/src/core/credentials/application"
credsDomain "github.com/gminetoma/go-auth/src/core/credentials/domain"
```

### Packages

- Lowercase, single word when possible: `domain`, `clock`, `errs`
- Never `util`, `helper`, `common`
- Don't stutter: `credentials.Credentials` is acceptable, `credentials.CredentialsService` is not

### Variables and Constants

```go
// Good
const maxRetries = 3
const DefaultTimeout = 30 * time.Second
var ErrNotFound = errors.New("not found")

// Bad
const MAX_RETRIES = 3
const DEFAULT_TIMEOUT = 30 * time.Second
```

### Initialisms

Keep consistent casing: `URL`, `ID`, `HTTP`, `SQL`, `API`.

```go
// Good
ownerID, userID, httpClient, sqlDB

// Bad
ownerId, userId, httpClient, sqlDb
```

### Receivers

Short, one or two letters, consistent across all methods of a type:

```go
func (c *Credentials) VerifyPassword(password string) bool
func (t *RefreshToken) IsExpired(now time.Time) bool
func (r *PGCredentialsRepository) Create(ctx context.Context, ...) error
```

### Getters

Do not use `Get` prefix:

```go
// Good
func GetOwnerID(ctx context.Context) (domain.OwnerID, bool)

// Bad — only when the underlying concept uses "get"
// Otherwise:
func OwnerID() domain.OwnerID
```

### Interfaces

- Define interfaces in the **consuming** package, not the implementing package
- Name single-method interfaces with `-er` suffix when natural: `Reader`, `Writer`

```go
// application/token_service.go — consumer defines the interface
type TokenService interface {
    Generate(ownerID domain.OwnerID) (string, error)
    Verify(token string) (domain.OwnerID, error)
}

// infrastructure/jwt_token_service.go — implementer doesn't reference the interface
```

## Context

- `context.Context` is always the **first parameter**
- Pass context through all I/O boundaries (DB, HTTP, external services)
- Never store context in a struct

```go
// Good
func (r *PGCredentialsRepository) Create(ctx context.Context, credentials *domain.Credentials) error

// Bad
func (r *PGCredentialsRepository) Create(credentials *domain.Credentials) error
```

## Error Handling

### Return pattern

`error` is always the **last return value**:

```go
func NewCredentials(params NewCredentialsParams) (*Credentials, error)
```

### Handle errors first

```go
// Good
creds, err := s.credentialsRepository.GetByEmail(ctx, email)
if err != nil {
    if errors.Is(err, errs.ErrNotFound) {
        return nil, ErrInvalidCredentials
    }

    return nil, err
}

// continue with creds
```

### Sentinel errors

Use unparameterized global values for distinguishable error conditions:

```go
var (
    ErrInvalidCredentials  = errors.New("auth.invalid-credentials")
    ErrInvalidRefreshToken = errors.New("auth.invalid-refresh-token")
)
```

### Error strings

- Do not capitalize (unless starting with a proper noun or exported name)
- Do not end with punctuation

```go
// Good
errors.New("auth.invalid-credentials")
fmt.Errorf("unexpected signing method: %v", alg)

// Bad
errors.New("Auth.Invalid-Credentials.")
```

### Error wrapping

- Use `%w` to preserve error chains for `errors.Is()` and `errors.As()`
- Use `%v` at system boundaries to avoid leaking internal errors

## Structs

### Initialization

Always use field names:

```go
// Good
return &Credentials{
    ID:           CredentialsID(id.Make()),
    OwnerID:      params.OwnerID,
    Email:        params.Email,
    PasswordHash: passwordHash,
}

// Bad
return &Credentials{
    CredentialsID(id.Make()),
    params.OwnerID,
    params.Email,
    passwordHash,
}
```

### Zero values

Omit zero-value fields when intent is clear:

```go
// OK — ReadAt is intentionally zero (unread notification)
return &Notification{
    ID:        NotificationID(id.Make()),
    OwnerID:   params.OwnerID,
    Title:     params.Title,
    CreatedAt: params.Now,
}
```

## Functions

### Option structs

Use parameter structs for functions with many parameters:

```go
type NewAuthServiceParams struct {
    CredentialsRepository  CredentialsRepository
    RefreshTokenRepository RefreshTokenRepository
    TokenService           TokenService
    RefreshTokenExpiry     time.Duration
    Clock                  clock.Clock
}

func NewAuthService(params NewAuthServiceParams) *AuthService
```

### Keep it simple

- Don't add error handling for impossible scenarios
- Don't create abstractions for one-time operations
- Don't design for hypothetical future requirements

## Testing

### Standard library

Use `testing` package. Avoid assertion libraries.

### Error messages

Include function name, inputs, actual result, and expected result:

```go
if got != want {
    t.Errorf("VerifyPassword(%q) = %v, want %v", password, got, want)
}
```

### Comparisons

Use `github.com/google/go-cmp/cmp` for deep equality instead of hand-coded field comparisons.

### Table-driven tests

Use field names in test case structs:

```go
tests := []struct {
    name     string
    email    string
    password string
    wantErr  bool
}{
    {name: "valid credentials", email: "admin@test.com", password: "secret", wantErr: false},
    {name: "empty email", email: "", password: "secret", wantErr: true},
}
```

### Helpers

Mark helpers with `t.Helper()`:

```go
func createTestCredentials(t *testing.T) *domain.Credentials {
    t.Helper()
    // ...
}
```

## Panics

- Never panic on transient failures — return errors
- Reserve panics for programmer errors (API misuse)
- Never panic in library code

## Comments

- Only comment non-obvious logic — don't restate what the code does
- Exported types and functions should have doc comments
- Start comments with the name of the thing being described

```go
// AuthService handles authentication operations including login,
// token refresh, and logout.
type AuthService struct { ... }
```
