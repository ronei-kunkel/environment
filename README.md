# Dynamic Environment Loader

[![Go Reference](https://pkg.go.dev/badge/github.com/ronei-kunkel/environment.svg)](https://pkg.go.dev/github.com/ronei-kunkel/environment)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue)](https://golang.org/doc/go1.21)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**Production-safe dynamic environment loader for Go** – loads environment variables into developer-defined structs and fails fast if any required variable is missing.

---

## 🔗 Module Dependencies

- [`github.com/joho/godotenv`](https://github.com/joho/godotenv)

---

## ⚙️ Behavior

- Maps environment variables from host or `.env` files into a Go struct.  
- Supports **custom variable names** via struct tags (`env`).  
- Aborts the program with `log.Fatalln` if any required field is missing.  
- Supports **multiple `.env` files**; later files override earlier ones.  
- Works with any **struct type** defined by the developer.  

---

## 📝 Usage

### 1. Define your struct

```go
// internal/env/vars.go
package env

type Vars struct {
  ENVIRONMENT string
  DB_NAME     string
}
```

### 2. Load environment variables

```go
// main.go
package main

import (
  "my-project/internal/env"

  "github.com/ronei-kunkel/environment"
)

func main() {
  envVars := environment.Load[env.Vars]()
}
```

### 3. Custom mapping with `env` tag

```go
type Vars struct {
  ENVIRONMENT string `env:"APP_ENV"`
  DB_NAME     string
  SOME_KEY    string
}
```

### 4. Define methods on your struct

```go
func (v Vars) IsProdEnv() bool {
  switch strings.ToUpper(v.ENVIRONMENT) {
  case "PROD", "PRODUCTION":
    return true
  }
  return false
}
```

### 5. Load custom `.env` files

```go
envVars := environment.Load[env.Vars](".prod.env")
envVars := environment.Load[env.Vars](".prod.env", ".fallback.env")
```

### 6. Error Handling

```txt
Errors loading environment variables:
 - has no `APP_ENV` environment variable defined to populate into `ENVIRONMENT`
 - has no `DB_NAME` environment variable defined to populate into `DB_NAME`
Aborting due to missing env vars
```

### 7. Best Practices

- Define **all required fields** in your struct.  
- Use `.env` files **only in development**.  
- Prefer **actual environment variables** in production.  
- Use `env:"..."` tags to map differently named variables.  
- Write **tests** to ensure missing variables abort as expected.  

### 8. Testing Suggestions

1. Use `os.Setenv()` to set environment variables during tests.  
2. For `log.Fatalln`, run `Load` in a **subprocess** to capture exit codes.  
3. Test multiple `.env` files and tag mappings.

### 9. Full Example

```go
type Vars struct {
  ENVIRONMENT string `env:"APP_ENV"`
  DB_NAME     string
  SOME_KEY    string
}

func main() {
  envVars := environment.Load[Vars](".prod.env", ".fallback.env")

  if envVars.IsProdEnv() {
    // Production-specific logic
  }
}
```

### 🔹 Tips

- Keep `.env` files **out of version control** if they contain secrets.  
- Use `.env` files **for local development only**.  
- This module is designed to **fail fast**, preventing runtime misconfiguration.
