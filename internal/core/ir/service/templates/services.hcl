services {

  # Global API configuration
  base_path   = "/api/v1"
  middlewares = ["log"]

  ################################################################
  # USER SERVICE
  ################################################################
  service "user" {

    model   = "user"
    expose  = true # Optional (defaults to true)
    prefix  = "/users"   # Optional override (defaults to /users (model name pluralized))

    ############################################################
    # Features: CRUD & Query Operators
    ############################################################
    features {
      pagination = true
      sorting    = ["id", "username", "email", "age", "isActive"]
      filtering  = ["email", "age", "isActive", "role", "organization"]
      search     = ["username", "email", "tags"]
    }

    ############################################################
    # CRUD Method Toggles
    ############################################################
    methods {
      create = true
      read   = true
      update = true
      delete = true
      list   = true
      bulkCreate = true # Enable bulk create endpoint by default false
      bulkUpdate = false # Enable bulk update endpoint by default false
      bulkDelete = false # Enable bulk delete endpoint by default false
    }

    ############################################################
    # Access Control (Role → Allowed Methods)
    ############################################################
    access {
      admin     = ["create", "read", "update", "delete", "list"]
      user      = ["read", "update"]
      moderator = ["read", "list"]

      # Guest users may sign up
      guest     = ["create"]
    }

    ############################################################
    # Optional service-level overrides
    ############################################################
    options {
      soft_delete       = true
      soft_delete_field = "deletedAt"

      rate_limit {
        per_minute = 200
        per_hour   = 5000
      }

      # Attach per-service middlewares
      middlewares = ["auth"]
    }
  }

  ################################################################
  # PRODUCT SERVICE
  ################################################################
  service "product" {

    model  = "product"
    expose = true
    prefix = "/products"

    # features
    pagination = true
    sorting    = ["id", "title", "price", "stock", "isActive"]
    filtering  = ["price", "isActive", "stock", "categories"]
    search     = ["title", "description", "categories"]

    methods {
      create = true
      read   = true
      update = true
      delete = true
      list   = true
    }

    access {
      admin     = ["create", "read", "update", "delete", "list"]
      user      = ["read", "list"]
      moderator = ["read", "list", "update"]
    }

    options {
      rate_limit {
        limit = "300/hour" # other options : (day, second, minute, hour, month, 10minutes, 10seconds etc.)
        action = "throttle"  # other options: "block", "log"
        # block or throttle for time after limit reached (default same as limit duration)
        # action_duration = "15m" # 15 minutes other options: (s, m, h)
        # key_by = "ip"  # other options: "user_id", "api_key", "header:<HEADER_NAME>"
      }
    }
  }

}
