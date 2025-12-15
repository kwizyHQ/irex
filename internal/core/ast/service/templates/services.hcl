

policies {
  
  mode = "deny-by-default" # default mode: deny-by-default, allow-by-default
  precedence = "deny-over-allow" # default precedence: deny-over-allow, allow-over-deny
  short_circuit = true # default short-circuit evaluation: true, false

  policy "authenticated" {
    effect = "allow"
    scope = "request"
    description = "Allows access only to authenticated users"
    rule = "ctx.auth != null"
  }

  policy "guest-readonly" {
    effect = "allow"
    scope = "request"
    description = "Allows read-only access to guests"
    rule = "ctx.auth == null ? ctx.request.method in ['GET','HEAD']"
  }

  policy "only-owner" {
    effect = "allow"
    scope = "resource"
    description = "Allows access only to resource owner"
    rule = "ctx.auth != null && ctx.auth.id == item.owner.id"
  }

  policy "admin-users" {
    effect = "allow"
    scope = "request"
    description = "Allows access to admin users"
    rule = "ctx.auth != null && ctx.auth.role == 'admin'"
  }

  custom "trial-users" {
    scope = "request"
    description = "Allows access to trial users"
    # custom policy to be implemented by user in respective file inside policies folder
  }

  group "write-access" {
    # NOTE: All policies in a group must have the same scope. Cross-scope grouping is not allowed.
    scope = "request"
    description = "Allows access to admin users"
    policies = ["admin-users","authenticated"]
  }
  
}

rate_limits {

  defaults {
    action = "throttle" # default action: throttle, block
    type = "fixed_window" # default type: fixed_window, sliding_window, token_bucket
    count_key = ["ctx.request.ip"] # default count key (array of strings)
    limit = "100/min" # default limit if not specified in preset or service level - only for fixed_window and sliding_window types
    # bucket_size = 500 # only for token_bucket type
    # refill_rate = "10/min" # only for token_bucket type
    # burst = 50 # only for token_bucket type
    response {
      status_code = 429
      body = {
        message = "Too many requests, please try again later."
      }
    }
  }
  preset "server-total" {
    limit = "10000/min" # 1000 requests per minute
    count_key = ["global"] # array of strings
  }
  preset "user-total" {
    limit = "500/min" # 500 requests per minute
    count_key = ["ctx.auth.id"] # array of strings
  }
  preset "ip-total" {
    limit = "200/min" # 200 requests per minute
  }
  preset "standard" {
    limit = "100/min" # 100 requests per minute
  }
  preset "request-method-limit" {
    limit = "200/min" # 200 requests per minute
    count_key = ["ctx.request.ip","ctx.request.method"] # array of strings
  }
  preset "standard-sliding" {
    type = "sliding_window"
    limit = "100/min"
  }
  preset "standard-token" {
    type = "token_bucket"
    refill_rate = "5/min" # refill 5 tokens per minute
    bucket_size  = 300
    burst = 50
  }

  custom "premium-users" {
    # custom rate limit preset to be implemented by user in respective file inside rate_limits folder
  }
}

services {

  # ---------------------------------------------------------------------------
  # Global API configuration
  # ---------------------------------------------------------------------------
  base_path = "/api/v1"

  cors = true
  allowed_origins = ["*"]
  allowed_methods = ["GET","POST","PUT","DELETE","PATCH","OPTIONS"]
  allowed_headers = ["Content-Type","Authorization"]
  expose_headers  = ["X-Total-Count"]
  allow_credentials = false
  max_age = 86400
  cache_control = "no-store"

  # ---------------------------------------------------------------------------
  # Global service defaults - All services will inherit
  # ---------------------------------------------------------------------------
  defaults {
    pagination = true
    expose = true
    soft_delete = true

    crud_operations = ["create","read","update","delete","list"]
    # batch_operations = ["batch_create","batch_update","batch_delete"]

    middlewares = ["log"]

    sorting   = ["created_at", "updated_at"]
    filtering = ["user_id","status","created_at"]
    # search = ["name","description"]
  }

  # ---------------------------------------------------------------------------
  # Global operations (non-model endpoints)
  # ---------------------------------------------------------------------------
  operation "health_check" {
    method = "GET"
    path   = "/health"
    description = "Health check endpoint"
  }

  # ---------------------------------------------------------------------------
  # User service
  # ---------------------------------------------------------------------------
  service "user" {
    model = "User"
    path  = "users"

    middlewares = ["auth","log"]

    apply "policy" "authenticated" {
      to_operations = ["*"]
      rate_limits = ["user-total"]
    }

    service "user-posts" {
      model = "Post"
      path  = "{user_id}/posts"

      crud_operations  = ["create","read","list"]
      batch_operations = ["batch_create"]
    }
  }

  # ---------------------------------------------------------------------------
  # Post service
  # ---------------------------------------------------------------------------
  service "post" {
    model = "Post"
    path  = "posts"

    # incorrect usage example of policy application at service level - Resource-scoped policies must not apply rate limits at all.
    # apply "policy" "only-owner" {
    #   rate_limits = ["standard"]
    # }

    apply "policy" "trial-users" {
      to_operations = ["create"]
      rate_limits = ["standard"]
    }

    apply "rate_limit" "premium-users" {
      to_operations = ["*"]
    }

    service "post-comments" {
      model = "Comment"
      path  = "{post_id}/comments"

      crud_operations = ["create","read","list"]

      operation "like" {
        method = "POST"
        path   = "{comment_id}/like"
        description = "Like a comment"
        action = "posts/increment_like_count"
      }

      operation "unlike" {
        method = "POST"
        path   = "{comment_id}/unlike"
        description = "Unlike a comment"
        action = "posts/decrement_like_count"
      }
    }
  }
}
