# goenv-loader

`goenv-loader` is a lightweight Go package for loading environment variables into your application’s configuration structs.  
It supports:  
✅ Automatic environment variable injection into struct fields  
✅ Default values when an environment variable is missing  
✅ Required field enforcement  
✅ Nested structs  
✅ Basic type conversion (`string`, `int`)  
✅ Built-in validation + extendable custom validation  

Perfect for building 12-factor apps in Go without wiring up boilerplate env parsing code.

---

## Installation

```bash
go get github.com/azpz30/goenv-loader
```

## Quick Start

Define your configuration struct with env, required, and default tags:

```go
package main

import (
	"fmt"
	"log"

	goenv "github.com/<your-username>/goenv-loader"
)

type Config struct {
	SystemID string `env:"SYSTEM_ID" required:"true"`
	LogLevel string `env:"LOG_LEVEL" default:"info"`
	ApiPort  string `env:"API_PORT" default:":8080"`
	DB       struct {
		Username string `env:"DATABASE_USERNAME"`
		Password string `env:"DATABASE_PASSWORD"`
		Port     int    `env:"DATABASE_PORT" default:"5432"`
	}
}

func main() {
	var cfg Config
	if err := goenv.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Printf("%+v\n", cfg)
}

```

## Example Environment

```bash
export SYSTEM_ID=rating-service
export DATABASE_USERNAME=admin
export DATABASE_PASSWORD=secret
```

### Output

```bash
{SystemID:rating-service LogLevel:info ApiPort::8080 DB:{Username:admin Password:secret Port:5432}}
```

## Features

- **Automatic environment loading**  
  Load environment variables directly into struct fields using reflection.

- **Required fields**  
  ```go
  SystemID string `env:"SYSTEM_ID" required:"true"`
  ```
  If SYSTEM_ID is missing, goenv-loader throws an error.

- **Default values** 
  ```go
  LogLevel string `env:"LOG_LEVEL" default:"info"`
  ```
  If LOG_LEVEL is not set, "info" will be used.

- **Nested structs** 
  Nested configurations are supported out-of-the-box, e.g.:
  ```go
  DB struct {
    Username string `env:"DATABASE_USERNAME"`
    Password string `env:"DATABASE_PASSWORD"`
  }
  ```
- **Type Conversion**
  Converts string environment variables into supported Go types:
  ```go
  string
  int
  ```
- **Basic validation**
  - Integers must be greater than 0
  - Strings must not be empty if required:"true"

- **Custom validation support**
  - Extend validation rules easily by wrapping your config in a higher-level package.
