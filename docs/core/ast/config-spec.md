
<h1 align="center" style="color:#10b981; font-size:2.4em; font-weight:bold;"> IREX Configuration Guide </h1>
<p align="center"> Define how your application is generated, structured, and run using a single, declarative configuration. </p>

## What is this configuration?

The IREX configuration file is a single source of truth that describes:

- What your project is
- What should be generated
- Which runtime and frameworks to target
- How schemas and services should behave

IREX uses this file to generate backend code consistently across supported runtimes and frameworks.

The configuration is written in HCL (HashiCorp Configuration Language) and is designed to be:

- Human-readable
- Declarative
- Easy to extend over time

---

## Overall Structure

At a high level, the configuration looks like this:

```hcl
project {
  paths { ... }
  generator { ... }
  runtime { ... }
  meta { ... }
}
```

Each block controls a specific aspect of your project.

---

## Project Information

This section describes your project for documentation and metadata purposes.

```hcl
project {
  name        = "MyApp"
  description = "An example backend generated with IREX"
  version     = "1.0.0"
  author      = "Your Team"
  license     = "MIT"
  timezone    = "UTC"
}
```

### Fields

| Field       | Type   | Description                      |
|-------------|--------|----------------------------------|
| name        | string | Human-friendly project name       |
| description | string | Short summary of the project      |
| version     | string | Project version                  |
| author      | string | Individual or organization        |
| license     | string | SPDX license identifier           |
| timezone    | string | Default timezone (e.g. UTC)       |

> ℹ️ These values do not change runtime behavior. They are informational.

---

## Paths

Paths tell IREX where to read inputs from and where to write outputs.

```hcl
paths {
  specifications = "./specs"
  templates      = "./templates"
  output         = "./generated"
}
```

### Fields

| Field          | Type   | Description                        |
|----------------|--------|------------------------------------|
| specifications | string | Directory containing IREX specifications |
| templates      | string | Template directory used by generators    |
| output         | string | Root directory for generated code       |

---

## Generator Settings

Controls what gets generated and how generation behaves.

```hcl
generator {
  schema       = true
  service      = true
  dry_run      = false
  clean_before = true
}
```

### Fields

| Field        | Type | Description                                 |
|--------------|------|---------------------------------------------|
| schema       | bool | Generate database / data schemas            |
| service      | bool | Generate API services                      |
| dry_run      | bool | Validate and build IR without writing files |
| clean_before | bool | Remove old generated files before generating|

---

## Runtime

Defines which runtime and frameworks your application targets.

```hcl
runtime {
  name     = "node-ts"
  version  = "18"
  scaffold = true
}
```

### Runtime Fields

| Field    | Type   | Allowed values (enum) | Description                |
|----------|--------|-----------------------|----------------------------|
| name     | string | node-ts, node-js      | Target runtime             |
| version  | string | (free-form)           | Target runtime version     |
| scaffold | bool   | true / false          | Generate base project structure |

---

## Runtime Options

Fine-grained runtime configuration.

```hcl
runtime {
  options {
    package_manager = "npm"
    entry           = "src/app.ts"
    dev_nodemon     = false
  }
}
```

### Fields

| Field           | Type   | Allowed values (enum) | Description                  |
|-----------------|--------|----------------------|------------------------------|
| package_manager | string | npm, yarn, pnpm      | Dependency manager           |
| entry           | string | —                    | Application entry file       |
| dev_nodemon     | bool   | true / false         | Enable hot reload in development |

---

## Schema Configuration

Controls data modeling and database layer.

```hcl
runtime {
  schema {
    framework = "mongoose"
    version   = "6"
    options {
      uri = env("DATABASE_URI")
      db  = env("DATABASE_NAME")
    }
  }
}
```

### Schema Fields

| Field     | Type   | Allowed values (enum) | Description           |
|-----------|--------|----------------------|-----------------------|
| framework | string | mongoose, knex       | Schema / ORM framework|
| version   | string | —                    | Framework version     |

#### Schema Options

| Field | Type        | Description                 |
|-------|-------------|----------------------------|
| uri   | env(string) | Database connection string  |
| db    | env(string) | Database name               |

> 🔐 Environment values are resolved at runtime.

---

## Service Configuration

Controls HTTP layer and API framework.

```hcl
runtime {
  service {
    framework = "fastify"
    version   = "4"
    options {
      logger = true
      port   = 8080
      host   = "localhost"
    }
  }
}
```

### Service Fields

| Field     | Type   | Allowed values (enum) | Description         |
|-----------|--------|----------------------|---------------------|
| framework | string | fastify, express     | HTTP framework      |
| version   | string | —                    | Framework version   |

#### Service Options

| Field | Type   | Description              |
|-------|--------|-------------------------|
| logger| bool   | Enable framework logging |
| port  | number | Server port              |
| host  | string | Server bind address      |

---

## Metadata

Optional information about when and how the configuration was created.

```hcl
meta {
  created_at = "2025-12-10"
  generator  = "irex"
}
```

### Fields

| Field      | Type   | Description           |
|------------|--------|----------------------|
| created_at | string | Creation date        |
| generator  | string | Generator identifier |

---

## Minimal Example

```hcl
project {
  name = "ExampleApp"

  paths {
    specifications = "./specs"
    output         = "./generated"
  }

  generator {
    schema  = true
    service = true
  }

  runtime {
    name     = "node-ts"
    version  = "18"
    scaffold = true

    options {
      package_manager = "npm"
      entry           = "src/app.ts"
    }

    schema {
      framework = "mongoose"
    }

    service {
      framework = "fastify"
      options {
        port = 8080
      }
    }
  }
}
```
