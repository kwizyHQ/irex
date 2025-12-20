
<h1 align="center" style="color:#f59e42; font-size:2.4em; font-weight:bold;"> IREX Schema & Model Guide </h1>
<p align="center"> Define your data models once, generate schemas and APIs everywhere. </p>

## What is a Schema in IREX?

The schema specification describes your application’s data models, including:

- Entities (models)
- Fields and validation rules
- Indexes
- Relationships
- Database-specific behavior

IREX treats schemas as runtime-agnostic definitions. The same schema can be used to generate models for different databases and frameworks.

---

## Overall Structure

All models live inside a single `models` block.

```hcl
models {
  model User { ... }
  model Post { ... }
  model Comment { ... }
}
```

Each model represents a real-world entity in your system.

---

## Model Definition

A model defines:

- Its fields
- Optional configuration
- Relationships to other models

```hcl
model User {
  field username { ... }
  field email { ... }

  config { ... }
  relations { ... }
}
```

### Model Properties

| Element   | Type    | Description                  |
|-----------|---------|------------------------------|
| model     | block   | Declares a new data model    |
| field     | block[] | Defines fields/properties    |
| config    | block   | Model-level configuration    |
| relations | block   | Relationships to other models|

---

## Field Definition

A field represents a single property on a model.

```hcl
field username {
  type     = "string"
  required = true
  unique   = true
  trim     = true
}
```

### Common Field Attributes

| Field      | Type   | Allowed values (enum)                | Description                |
|------------|--------|--------------------------------------|----------------------------|
| type       | string | string, number, boolean, date, object, array, enum | Data type                  |
| required   | bool   | true / false                         | Is the field mandatory     |
| unique     | bool   | true / false                         | Must value be unique       |
| default    | any    | —                                    | Default value              |
| description| string | —                                    | Human-readable description |
| visibility | string | public, private, internal            | Exposure level             |

---

## Validation Rules

Validation rules depend on the field type.

### String Fields

| Field     | Type   | Description         |
|-----------|--------|--------------------|
| minlength | number | Minimum length      |
| maxlength | number | Maximum length      |
| trim      | bool   | Trim whitespace     |
| match     | string | Regex pattern       |
| message   | string | Custom error message|

### Number Fields

| Field | Type   | Description    |
|-------|--------|---------------|
| min   | number | Minimum value  |
| max   | number | Maximum value  |

---

## Nested / Embedded Fields

Fields can contain sub-fields to represent embedded objects.

```hcl
field profile {
  type = "object"

  fields {
    field age {
      type = "number"
    }

    field bio {
      type = "string"
    }
  }
}
```

This is useful for JSON objects or embedded documents.

---

## Database-Specific Field Options

Field-level database customization is optional and non-portable.

### MongoDB Field Options

```hcl
mongo {
  index  = true
  unique = true
}
```

| Field | Type | Description         |
|-------|------|--------------------|
| index | bool | Create index        |
| unique| bool | Enforce uniqueness  |

### SQL / MySQL Field Options

```hcl
mysql {
  index   = true
  unique  = true
  collate = "utf8mb4_general_ci"
}
```

| Field   | Type   | Description         |
|---------|--------|--------------------|
| index   | bool   | Create index        |
| unique  | bool   | Enforce uniqueness  |
| collate | string | Column collation    |

---

## Model Configuration

Controls how the model behaves as a whole.

```hcl
config {
  timestamps = true
  table      = "users"
}
```

### Model Config Fields

| Field      | Type   | Allowed values (enum) | Description                      |
|------------|--------|----------------------|----------------------------------|
| timestamps | bool   | true / false         | Add created/updated timestamps   |
| table      | string | —                    | Custom table / collection name   |
| strict     | bool   | true / false         | Enforce schema strictly          |
| idStrategy | string | auto, uuid, custom   | ID generation strategy           |
| description| string | —                    | Model description                |

---

## Index Definitions

Indexes improve query performance and enforce constraints.

```hcl
index "user_email_idx" {
  fields = ["email"]
  unique = true
}
```

### Index Fields

| Field  | Type     | Description              |
|--------|----------|-------------------------|
| fields | string[] | Fields included in index |
| unique | bool     | Enforce uniqueness       |

---

## Database-Specific Model Options

### MongoDB Model Options

```hcl
mongo {
  collection  = "users"
  autoIndex   = true
  strictQuery = true
}
```

| Field       | Type   | Description              |
|-------------|--------|-------------------------|
| collection  | string | Collection name         |
| autoIndex   | bool   | Auto-create indexes     |
| autoCreate  | bool   | Auto-create collection  |
| strictQuery | bool   | Enforce strict queries  |

### SQL / MySQL Model Options

```hcl
mysql {
  engine  = "InnoDB"
  collate = "utf8mb4_general_ci"
}
```

| Field   | Type   | Allowed values (enum) | Description         |
|---------|--------|----------------------|--------------------|
| engine  | string | InnoDB, MyISAM       | Storage engine     |
| collate | string | —                    | Table collation    |

---

## Relationships

Relationships define how models connect to each other.

```hcl
relations {
  hasMany posts {
    ref = "Post"
  }

  belongsTo user {
    ref      = "User"
    onDelete = "CASCADE"
  }
}
```

### Supported Relation Types

| Relation    | Description         |
|-------------|--------------------|
| hasOne      | One-to-one         |
| hasMany     | One-to-many        |
| belongsTo   | Inverse relation   |
| manyToMany  | Many-to-many       |

### Common Relation Fields

| Field    | Type   | Allowed values (enum) | Description         |
|----------|--------|----------------------|--------------------|
| ref      | string | —                    | Target model       |
| onDelete | string | CASCADE, RESTRICT, SET_NULL | Delete behavior   |
| onUpdate | string | CASCADE, RESTRICT    | Update behavior    |

---

## Minimal Example

```hcl
models {
  model Example {
    field name {
      type     = "string"
      required = true
    }

    config {
      timestamps = true
    }
  }
}
```

---

## Design Principles

- Declarative, not database-first
- Portable across runtimes
- Database-specific options are optional
- Safe defaults with explicit overrides
- Relations are semantic, not implementation-bound
