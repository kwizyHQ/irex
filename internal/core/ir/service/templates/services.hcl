

policies {
  defaults {
    deny_rule = true
  }
  preset "authenticated" {
    description = "Allows access to authenticated users"
    allow_rule = "ctx.auth != null"
    deny_rule = true # optional override to deny by default unless allow_rule matches
  }
  preset "only-owner" {
    description = "Allows access only to the owner of the resource"
    allow_rule = "ctx.auth.id == item.owner.id"
  }
  preset "admin-access" {
    description = "Allows access to admin users"
    allow_rule = "ctx.auth.role == 'admin'"
  }
  
  preset "guest-acess" {
    description = "Allows access to guest users"
    allow_rule = "ctx.auth == null"
  }
  preset "guest-readonly" {
    description = "Allows read-only access to guest users"
    allow_rule = "ctx.auth == null && (ctx.request.method 'GET' || ctx.request.method 'HEAD')"
  }
  # for custom policies, generator should leave empty and let user implement the logic in respective file 
  # inside the policies folder with name custom_{policy_name}.(extension_any supported_by_framework)
  preset "custom.allow-trial-users" {}

  # Policy mixins to compose complex policies from simpler ones these group policies can be used in services just like normal policies
  group "admin-or-owner" {
    description = "Allows access to admin users or resource owners"
    any = ["admin-access", "only-owner"] # mix of two policies definations 
    all = [] # all policies in this list must pass
  }
}

rate_limits {

  defaults {
    action = "throttle" # default action: throttle, block
    type = "fixed_window" # default type: fixed_window, sliding_window, token_bucket
    count_key = "ctx.request.ip" # default count key
    # limit = "1000/hour" # default limit if not specified in preset or service level - only for fixed_window and sliding_window types
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
  preset "standard" {
    limit = "100/min" # 100 requests per minute
  }
  preset "get-heavy" {
    limit = "500/min"
    expression = "ctx.request.method == 'GET'" # only for GET requests
  }
  preset "sliding-100-per-minute" {
    type = "sliding_window"
    limit = "100/min"
  }
  preset "token-bucket-300-per-hour" {
    type = "token_bucket"
    refill_rate = "5/min" # refill 5 tokens per minute
    bucket_size  = 300
  }

  preset "custom.premium-users" {
    # custom rate limit preset to be implemented by user in respective file inside rate_limits folder
  }
}

services {

  # Global services configuration like base_path, 
  base_path = "/api/v1"
  cors = true # enable CORS for all services, default false
  allowed_origins = ["*"] # CORS allowed origins, default "*"
  allowed_methods = ["GET","POST","PUT","DELETE","PATCH","OPTIONS"] # CORS allowed methods, default all common methods
  allowed_headers = ["Content-Type","Authorization"] # CORS allowed headers, default common headers
  expose_headers = ["X-Total-Count"] # CORS expose headers, default none
  allow_credentials = false # CORS allow credentials, default false
  max_age = 86400 # CORS max age in seconds, default 86400 (1 day)
  cache_control = "no-store" # default cache control header for all services, default "no-store"


  
  # Service-level defaults applied to all services unless overridden
  defaults {
    pagination = true # defaults to false - can be set to enable pagination by default
    expose = true # defaults to true - can be set so user need explicitly expose a service
    soft_delete = true # defaults to false - can be set to enable soft delete by default (needs model support)
    crud_operations    = ["create","read","update","delete","list"] # defaults to all CRUD methods plus list
    # batch_operations = ["batch_create","batch_update","batch_delete"] # defaults to batch variants of create, update, delete
    middlewares = ["log"] # default middlewares applied to all services (must exists in middlewares folder to be valid)
    sorting = ["created_at", "updated_at"] # default sorting fields - can be set to enable sorting by default
    filtering = ["user_id","status","created_at"] # default filtering fields - can be set to enable filtering by default
    # search = ["name","description"] # default search fields - can be set to enable search by default
  }

  /* Global rate limiting configuration
  * Applies to all services, total rate limit will not exceed these values unless overridden at service level
  * It may contain multiple rate limiters with different time windows and settings AND OR supported
  */

  operation "health_check" {
    # method = "GET" # defaults to GET - optional
    path   = "/health"
    description = "Health check endpoint"
    # action = "HealthController.check" # optional handler to wire to controller stub without handler request will return 200 OK
  }

  service "user" {
    model = "User" # references models.hcl - this is the main data model will be used for crud operations
    expose = true  # whether to expose this service via generated API routes - default true
    path = "users" # relative to base_path - defaults to plural of service name
    # crud_operations = ["create","read","update","delete"] # overrides defaults
    middlewares = ["auth","log"] # additional middlewares for this service
    policies = ["authenticated"] # policies to apply to this service
    service "user-posts" {
      model = "Post"
      path = "{user_id}/posts"
      crud_operations = ["create","read","list"]
      batch_operations = ["batch_create"]
      policies = ["only-owner"] # service-level policies can also be defined here
    }
  }

  service "post" {
    model = "Post"
    expose = true
    path = "posts"
    policies = ["guest-readonly","authenticated"]
    rate_limit {
      limit = "10/min" # overrides global limit for this service
      action = "throttle"
      action_duration = "1m"
      count_key = "ctx.auth.id"
    }
    service "post-comments" {
      model = "Comment"
      path = "{post_id}/comments"
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