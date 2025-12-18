models {

  model User {
    field username {
      type      = "string"
      required  = true
      unique    = true
      trim      = true
      minlength = 3
      maxlength = 30
      db {
        # database specific configurations can be added here that cannot be generalized above in fields settings
        mongo {
          collation {
            locale   = "en"
            strength = 1
          } # example
        }
        mysql {
          collate = "utf8mb4_general_ci"
        }
      }
    }

    field email {
      type     = "string"
      required = true
      unique   = true
      match    = "/^\\S+@\\S+\\.\\S+$/"
    }

    field password {
      type       = "string"
      required   = true
      minlength  = 8
      visibility = "private" # other options: "public" - (default) no need to define | "protected" - for using in models can be sent in responses with additional checks | "private" - used within the model only cannot override to sent in responses | "internal" (experimental) - cannot be used anywhere in the API only direct db access
    }

    field age {
      type     = "number"
      required = false
      default  = 18
      min      = 0
      max      = 150
    }

    field is_active {
      type     = "boolean"
      required = false
      default  = true
    }

    field organization {
      type     = "ref#organization"
      required = false
    }

    field tags {
      type     = "string[]"
      required = false
      default  = []
    }

    field phone {
      type     = "string"
      required = false
      match    = "^\\d{10}$"
      message  = "Phone must be a valid 10-digit number"
    }

    field address {
        field street {
          type      = "string"
          required  = true
          trim      = true
          maxlength = 100
        }
        field city {
          type      = "string"
          required  = true
          trim      = true
          maxlength = 50
        }
        field zip_code {
          type     = "number"
          required = true
          match    = "/^\\d{5}$/"
        }
        field country {
          type     = "string"
          required = true
          default  = "USA"
        }
    }

    config {
      timestamps = true
      table      = "app_users"
      strict     = true
      index "username_email_idx" {
        fields = ["username", "email"]
        unique = true
      }
      index "age_idx" {
        fields = ["age"]
        unique = false
      }
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
      manyToMany roles {
        ref = "Role"
      }
      hasMany posts {
        ref = "Post"
      }
    }
  }

  model Role {
    field name {
      type      = "string"
      required  = true
      unique    = true
      trim      = true
      minlength = 3
      maxlength = 30
    }

    field permissions {
      type     = "string[]"
      required = false
      default  = []
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
      manyToMany users {
        ref = "User"
      }
    }
  }

  model Post {
    field title {
      type        = "string"
      required    = true
      max       = 200
      description = "Post title"
    }

    field content {
      type        = "string"
      required    = false
      description = "Post content (body, markdown or HTML)"
    }

    field published {
      type        = "boolean"
      default     = false
      description = "Publication state"
    }

    config {
      table       = "posts"
      idStrategy  = "uuid"
      description = "User-generated posts"
      timestamps  = true
    }

    relations {
      belongsTo user {
        ref      = "User"
        onDelete = "CASCADE"
        onUpdate = "CASCADE"
      }
    }
  }

  model Comment {
    field content {
      type        = "string"
      required    = true
      description = "Comment content"
    }

    config {
      table       = "comments"
      idStrategy  = "uuid"
      description = "Comments on posts"
      timestamps  = true
    }

    relations {
      belongsTo post {
        ref      = "Post"
        onDelete = "CASCADE"
        onUpdate = "CASCADE"
      }
      belongsTo user {
        ref      = "User"
        onDelete = "SET NULL"
        onUpdate = "CASCADE"
      }
    }
  }
}