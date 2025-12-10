# =========================================================
# Irex Internal Representation - Schema Configuration
# =========================================================
# This file defines data models, their fields, relationships,
# and database-specific configurations for different supported databases.
#
# The specification below documents the supported configuration
# options that the IREX IR and generators understand. Options are
# split into general (cross-DB) settings and database-specific
# overrides under `db { mongo { ... } mysql { ... } }` or per-field
# `db { mongo { ... } sql { ... } }` blocks where needed.
#
# Guidance:
# - Keep cross-DB semantics at the top-level (for portability).
# - Place DB-specific settings inside a `db` block to avoid leaking
#   engine-specific behavior into the generic model representation.
# - Use `table` as the neutral term for storage naming (maps to
#   `collection` for MongoDB and `table` for SQL databases).
#
# =========================================================
# Config Specifications (top-level inside each `config` block)
# =========================================================
# Common (cross-DB) options
# - `table` (string): Neutral name for storage (collection/table).
# - `timestamps` (bool): Auto-manage createdAt/updatedAt fields.
# - `idStrategy` (string): `auto`, `uuid`, `snowflake`, `custom`.
# - `description` (string): Human-friendly model description.
# - `strict` (bool): Enforce schema validation on writes.
# - `softDelete` (bool): Use soft-delete (e.g. `deletedAt`) instead of hard delete.
# - `indexes` (list): Index definitions (see Index section).
# - `hooks` (object): Hook file paths or named hooks to attach.
# - `permissions` (object): Optional default ACLs for the model.
#
# Database-specific (`db` block inside `config`)
# - `db.mongo` (object): MongoDB-specific configuration examples:
#     - `collection` (string): Actual MongoDB collection name.
#     - `versionKey` (bool|string): Keep/rename `__v` key.
#     - `toJSON_getters` (bool): Enable toJSON getters behavior.
#     - `minimize` (bool): Mongoose minimize option for empty objects.
#     - `autoIndex` (bool): Enable auto-indexing in mongoose.
#     - `autoCreate` (bool): Auto-create collections on startup.
#     - `strictQuery` (bool): Strict query parsing option.
#     - `collation` (object): Default collection collation settings.
# - `db.mysql` or `db.sql` (object): SQL-specific configuration examples:
#     - `engine` (string): Table engine (e.g. `InnoDB`).
#     - `charset` (string): Character set (e.g. `utf8mb4`).
#     - `collation` (string): Collation (e.g. `utf8mb4_general_ci`).
#     - `tableName` (string): Explicit SQL table name override.
#     - `autoIncrement` (bool): Use auto-increment for numeric PKs.
#
# Example `config` block:
# config {
#   table = "users"
#   timestamps = true
#   idStrategy = "uuid"
#   strict = true
#   indexes = [ { fields = ["email"], unique = true } ]
#   db {
#     mongo = {
#       collection = "users"
#       versionKey = false
#     }
#     mysql = {
#       engine = "InnoDB"
#       charset = "utf8mb4"
#     }
#   }
# }
#
# =========================================================
# Index Specifications (inside `config.indexes`)
# =========================================================
# Each index entry is an object with these fields:
# - `fields` (list[string]): Ordered list of field names or nested paths.
# - `unique` (bool): Make the index unique.
# - `sparse` (bool): Sparse index behavior (Mongo specific).
# - `name` (string): Optional explicit index name.
# - `type` (string): Index type (e.g. `btree`, `hash`, `text`, `gin`).
# - `options` (object): DB-specific options (e.g., collation).
#
# Example indexes:
# indexes = [
#   { fields = ["username"], unique = true },
#   { fields = ["email"], unique = true, name = "ix_user_email" },
#   { fields = ["age"] }
# ]
#
# For DB-specific tuning you may put `db` key inside index:
# { fields = ["title"], type = "text", db = { mongo = { weights = { title = 10 } } } }
#
# =========================================================
# Field Specifications (inside `fields` blocks)
# =========================================================
# Supported field-level options (cross-DB first):
# - `type` (string): Basic data type or reference. Examples:
#     - Primitives: `string`, `number`, `boolean`, `date`, `any`
#     - Arrays: `string[]`, `number[]`, `ref#Model[]`
#     - References: `ref#Model` (for fk-like relations)
#     - Enum: `enum` plus `values` array
# - `required` / `optional` (bool): Presence constraint.
# - `unique` (bool): Uniqueness for the field.
# - `default` (any): Default value.
# - `min` / `max` (number): Numeric bounds.
# - `minlength` / `maxlength` (number): String length bounds.
# - `match` (string/regex): Regex validation.
# - `trim` (bool): Trim strings.
# - `visibility` (string): `public` | `protected` | `private` | `internal`.
# - `description` (string): Field description used in docs.
# - `fields` (object): Nested fields for object types.
# - `db` (object): DB-specific overrides for the field (see below).
#
# DB-specific field options (examples):
# - `db.mongo`:
#     - `type` (string): Explicit Mongo/BSON type if required.
#     - `collation` (object): Per-field collation for indexes.
#     - `index` (bool/object): Quick index hint or full index object.
# - `db.sql` / `db.mysql`:
#     - `type` (string): Column SQL type (e.g., `varchar(255)`, `tinyint(1)`).
#     - `nullable` (bool): Allow NULLs.
#     - `default` (value): Column default specific to SQL dialect.
#     - `extra` (string): Extra column args (e.g., `AUTO_INCREMENT`).
#
# Example field entries:
# username {
#   type = "string"
#   required = true
#   unique = true
#   minlength = 3
#   maxlength = 30
#   db = {
#     mongo = { collation = { locale = "en", strength = 2 } }
#     mysql = { collation = "utf8mb4_general_ci", type = "varchar(50)" }
#   }
# }
#
# Example cross-type override (boolean stored as tinyint in SQL):
# published {
#   type = "boolean"
#   default = false
#   db = { sql = { type = "tinyint(1)" }, mongo = { type = "boolean" } }
# }
#
# =========================================================
# Relation Specifications (inside `relations` blocks)
# =========================================================
# Relation object fields:
# - `ref` (string): Target model name (case-sensitive per IR).
# - `type` (string): Relation type: `hasOne`, `hasMany`, `belongsTo`, `manyToMany`.
# - `localField` / `foreignField` (string): Field names when explicit.
# - `through` (string): Join/through table or linking collection (for manyToMany).
# - `onDelete` / `onUpdate` (string): `CASCADE`, `SET NULL`, `RESTRICT`, etc.
# - `embedded` (bool): For document DBs, whether relation is embedded.
#
# Examples:
# # Simple belongsTo
# user {
#   ref = "User"
#   type = "belongsTo"
#   onDelete = "CASCADE"
# }
#
# # Many-to-many using a join table
# roles {
#   ref = "Role"
#   type = "manyToMany"
#   through = "user_roles"
# }
#
# # Embedded relation for Mongo
# meta {
#   ref = "Meta"
#   type = "hasOne"
#   embedded = true
# }
#
# =========================================================
# Notes and best practices
# =========================================================
# - Prefer defining behavior at a cross-DB level where possible.
# - Use `db` overrides only when engine-specific behavior is necessary.
# - Keep naming consistent: use `table` top-level, and map to
#   `collection` or `tableName` in DB-specific blocks.
# - Document any custom `idStrategy` or hook usage in the model's `description`.
#



models {

  User {
    fields {
      username {
        type      = "string"
        required  = true
        unique    = true
        trim      = true
        minlength = 3
        maxlength = 30
        db {
          # database specific configurations can be added here that cannot be generalized above in fields settings
          mongo {
            collation = {
              locale   = "en"
              strength = 1
            } # example
          }
          mysql {
            collation = "utf8mb4_general_ci"
          }
        }
      }

      email {
        type     = "string"
        required = true
        unique   = true
        match    = "/^\\S+@\\S+\\.\\S+$/"
      }

      password {
        type       = "string"
        required   = true
        minlength  = 8
        visibility = "private" # other options: "public" - (default) no need to define | "protected" - for using in models can be sent in responses with additional checks | "private" - used within the model only cannot override to sent in responses | "internal" (experimental) - cannot be used anywhere in the API only direct db access
      }

      age {
        type     = "number"
        required = false
        default  = 18
        min      = 0
        max      = 150
      }

      is_active {
        type     = "boolean"
        required = false
        default  = true
      }

      organization {
        type     = "ref#organization"
        required = false
      }

      tags {
        type     = "string[]"
        required = false
        default  = []
      }

      phone {
        type     = "string"
        required = false
        match    = "^\\d{10}$"
        message  = "Phone must be a valid 10-digit number"
      }

      address {
        fields {

          street {
            type      = "string"
            required  = true
            trim      = true
            maxlength = 100
          }
          city {
            type      = "string"
            required  = true
            trim      = true
            maxlength = 50
          }
          zip_code {
            type     = "number"
            required = true
            match    = "/^\\d{5}$/"
          }
          country {
            type     = "string"
            required = true
            default  = "USA"
          }
        }
      }
    }

    config {
      timestamps = true
      table      = "app_users"
      strict     = true
      indexes = [
        { fields = ["username"], unique = true },
        { fields = ["email"], unique = true },
        { fields = ["age"] }
      ]
      db {
        mongo {
          versionKey     = false
          collection     = "app_users"
          toJSON_getters = true
          minimize       = true
          autoIndex      = true
          autoCreate     = true
          strictQuery    = true
        }

        mysql {
          # MySQL-specific configuration placeholders
        }
      }
    }

    relations {
      roles {
        ref  = "Role"
        type = "manyToMany"
      }
    }
  }

  Role {
    fields {
      name {
        type      = "string"
        required  = true
        unique    = true
        trim      = true
        minlength = 3
        maxlength = 30
      }

      permissions {
        type     = "string[]"
        required = false
        default  = []
      }
    }

    config {
      timestamps = true
      table      = "roles"
      strict     = true

      db {
        mongo {
          versionKey     = false
          toJSON_getters = true
          minimize       = true
          autoIndex      = true
          autoCreate     = true
          strictQuery    = true
        }

        mysql {
          # MySQL-specific configuration placeholders
        }
      }
    }
    relations {
      users {
        ref  = "User"
        type = "manyToMany"
      }
    }
  }

  Post {
    config {
      table       = "posts"
      idStrategy  = "uuid"
      description = "User-generated posts"
      timestamps  = true
    }

    fields {

      # title and content
      title {
        type        = "string"
        required    = true
        limit       = 200
        description = "Post title"
      }

      content {
        type        = "string"
        optional    = true
        description = "Post content (body, markdown or HTML)"
      }

      published {
        type        = "boolean"
        default     = false
        description = "Publication state"
        db {
          sql {
            type = "tinyint(1)"
          }
          mongo {
            type = "boolean"
          }
        }
      }

    }
    relations {
      # Relationship: many posts belong to one user
      user {
        ref      = "User"
        type     = "belongsTo"
        onDelete = "CASCADE"
        onUpdate = "CASCADE"
      }
    }
  }

  Comment {
    config {
      table       = "comments"
      idStrategy  = "uuid"
      description = "Comments on posts"
      timestamps  = true
    }

    fields {
      content {
        type        = "string"
        required    = true
        description = "Comment content"
      }
    }

    relations {
      # Relationship: many comments belong to one post
      post {
        ref      = "Post"
        type     = "belongsTo"
        onDelete = "CASCADE"
        onUpdate = "CASCADE"
      }

      # Relationship: many comments belong to one user
      user {
        ref      = "User"
        type     = "belongsTo"
        onDelete = "SET NULL"
        onUpdate = "CASCADE"
      }
    }
  }
}