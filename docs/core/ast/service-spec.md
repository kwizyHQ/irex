


<h1 align="center" style="color:#3b82f6; font-size:2.5em; font-weight:bold; margin-bottom:0.5em;">Service Specifications</h1>

# Rate Limiting System

The IREX rate limiting system provides flexible, declarative controls for API usage. Rate limits are defined in the IR and can be applied globally, per-service, per-operation, or as custom presets. The system supports multiple strategies and is designed to be explicit and generator-friendly.

## Core Concepts

- **Explicit application**: No rate limit is active unless attached via an `apply` block.
- **Preset and custom limits**: Use built-in presets or define your own.
- **Multiple strategies**: Supports fixed window, sliding window, and token bucket algorithms.
- **Request-scoped only**: Rate limits are always evaluated before resource loading and must be request-scoped.

## Rate Limit Configuration

Rate limits are defined in the `rate_limits` block:

```hcl
rate_limits {
  defaults {
    action = "throttle" # throttle or block
    type = "fixed_window" # fixed_window, sliding_window, token_bucket
    count_key = ["ctx.request.ip"] # array of strings
    limit = "100/min"
    response {
      status_code = 429
      body = {
        message = "Too many requests, please try again later."
      }
    }
  }

  preset "server-total" {
    limit = "10000/min"
    count_key = ["global"] # array of strings
  }
  preset "user-total" {
    limit = "500/min"
    count_key = ["ctx.auth.id"] # array of strings
  }
  preset "ip-total" {
    limit = "200/min"
    count_key = ["ctx.request.ip"] # array of strings
  }
  preset "standard" {
    limit = "100/min"
  }
  preset "request-method-limit" {
    limit = "200/min"
    count_key = ["ctx.request.ip", "ctx.request.method"] # array of strings
  }
  preset "standard-sliding" {
    type = "sliding_window"
    limit = "100/min"
  }
  preset "standard-token" {
    type = "token_bucket"
    refill_rate = "5/min"
    bucket_size  = 300
    burst = 50
  }

  custom "premium-users" {
    # Custom rate limit preset to be implemented by user
  }
}
```

### Fields

- `action`: What to do when the limit is hit (`throttle` or `block`).
- `type`: Algorithm used (`fixed_window`, `sliding_window`, `token_bucket`).
- `count_key`: Array of strings. What to count against (IP, user ID, etc.). Always use an array, even for a single value (e.g., `["ctx.request.ip"]`).
- `limit`: The allowed number of requests per time window.
- `response`: Custom response for limited requests.
- `refill_rate`, `bucket_size`, `burst`: Token bucket-specific fields.

## Applying Rate Limits

Rate limits are attached using the `apply` block, either directly or via a policy:

```hcl
# Direct application
apply "rate_limit" "premium-users" {
	to_operations = ["*"]
}

# With a policy
apply "policy" "trial-users" {
	to_operations = ["create"]
	rate_limits = ["standard"]
}
```

### Rules

- Rate limits are always request-scoped and evaluated before resource loading.
- Resource-scoped policies must never apply rate limits.
- The referenced rate limit preset or custom must exist in the registry.
- `apply "rate_limit"` is unconditional; no policy or condition is involved.

## Example Usage

```hcl
service "post" {
	model = "Post"
	path  = "posts"

	apply "policy" "trial-users" {
		to_operations = ["create"]
		rate_limits = ["standard"]
	}

	apply "rate_limit" "premium-users" {
		to_operations = ["*"]
	}
}
```

## Invalid Usage

```hcl
# ❌ Resource-scoped policies must not apply rate limits
apply "policy" "only-owner" {
	rate_limits = ["standard"]
}
```


# Services System

The IREX services block defines the API surface of an application.
It is responsible for describing routes, operations, middleware, authorization hooks, and rate limiting behavior in a declarative, framework-agnostic way.

The services system is intentionally explicit:
nothing is applied automatically — all behavior must be attached using apply blocks.

## Core Design Principles

- **Explicit over implicit**
	- No policies or rate limits are active unless explicitly applied.
- **Separation of definition and activation**
	- Policies and rate limits are defined elsewhere, but activated only via apply.
- **Generator-friendly**
	- No merging, precedence resolution, or hidden inheritance logic is required.
- **Framework-agnostic**
	- Services describe what should happen, not how a framework implements it.

## Global API Configuration

The services block may define global API-level configuration:

```hcl
services {
	base_path = "/api/v1"

	cors = true
	allowed_origins = ["*"]
	allowed_methods = ["GET","POST","PUT","DELETE","PATCH","OPTIONS"]
	allowed_headers = ["Content-Type","Authorization"]
	expose_headers  = ["X-Total-Count"]
	allow_credentials = false
	max_age = 86400
	cache_control = "no-store"
}
```

These settings apply to all generated routes.

## Service Defaults

Global defaults define common behavior inherited by all services unless overridden.

```hcl
defaults {
	pagination = true
	expose = true
	soft_delete = true

	crud_operations = ["create","read","update","delete","list"]
	middlewares = ["log"]

	sorting   = ["created_at", "updated_at"]
	filtering = ["user_id","status","created_at"]
}
```

Defaults do not activate policies or rate limits.

## Services

A service represents a logical API resource, usually backed by a model.

```hcl
service "user" {
	model = "User"
	path  = "users"
}
```

Services may be nested to represent hierarchical routes.

## Operations

Operations define non-CRUD endpoints or custom actions.

```hcl
operation "health_check" {
	method = "GET"
	path   = "/health"
	description = "Health check endpoint"
}
```

Operations can exist globally or inside a service.

## The apply Block (Central Concept)

The apply block is the only mechanism for attaching behavior to services or operations.

Definitions are inert by default.
Nothing runs unless explicitly applied.

### Applying Policies

Policies are attached using `apply "policy"`.

```hcl
apply "policy" "authenticated" {
	to_operations = ["*"]
}
```

#### Rules

- The referenced policy must exist in the policy registry
- The policy must be request-scoped
- Resource-scoped policies cannot be used with rate limits
- `to_operations = ["*"]` means all operations in that service

### Applying a Policy with Rate Limits

Policies can optionally attach rate limits:

```hcl
apply "policy" "trial-users" {
	to_operations = ["create"]
	rate_limits   = ["standard"]
}
```

#### Important Notes

- Rate limits are evaluated before resource loading
- Rate limits are request-scoped only
- Resource-scoped policies must never apply rate limits

### Applying Rate Limits Directly

Rate limits can also be applied unconditionally, without a policy trigger.

```hcl
apply "rate_limit" "premium-users" {
	to_operations = ["*"]
}
```

This is equivalent to saying:

“Always apply this rate limit to these operations.”

#### Rules

- `apply "rate_limit"` is unconditional
- No policy or condition is involved
- The rate limit preset must exist

## Nested Services Example

```hcl
service "post" {
	model = "Post"
	path  = "posts"

	apply "policy" "authenticated" {
		to_operations = ["*"]
	}

	service "post-comments" {
		model = "Comment"
		path  = "{post_id}/comments"

		crud_operations = ["create","read","list"]
	}
}
```

Nested services inherit:

- base path
- defaults
- parent service context

## Invalid Usage (Enforced by Validator)

The following patterns are not allowed:

```hcl
# ❌ Resource-scoped policies must not apply rate limits
apply "policy" "only-owner" {
	rate_limits = ["standard"]
}

# ❌ No implicit policy or rate limit activation
service "post" {
	policies = ["authenticated"]     # deprecated
	rate_limits = ["standard"]       # deprecated
}
```

## Policies System

The IREX service layer supports a flexible, declarative policies system for authorization and access control. Policies are defined in the service IR (see `services.hcl`) and can be applied globally, per-service, per-operation, or grouped for reuse. The system is designed to be expressive, composable, and easy to extend.

### Policy Modes

- **mode**: Determines the default access behavior.
	- `deny-by-default`: All incoming requests are denied unless explicitly allowed by a policy.
	- `allow-by-default`: All incoming requests are allowed unless explicitly denied by a policy.

### Policy Precedence

- **precedence**: Controls how conflicting allow/deny policies are resolved.
	- `deny-over-allow`: Deny policies take priority over allow policies.
	- `allow-over-deny`: Allow policies take priority over deny policies.

### Short-Circuit Evaluation

- **short_circuit**: If `true`, policy evaluation stops as soon as a definitive allow/deny is determined.

### Policy Definition

Each policy must specify:

- `effect`: `allow` or `deny` (required)
- `scope`: `request` or `resource` (required)
- `rule`: Expression evaluated against the request context (required)
- `description`: (optional) Human-readable explanation

#### Example Preset Policies

- **authenticated**: Allows access only to authenticated users
	- `rule = "ctx.auth != null"`
- **guest-readonly**: Allows read-only access to guests
	- `rule = "ctx.auth == null ? ctx.request.method in ['GET','HEAD']"`
- **only-owner**: Allows access only to the resource owner
	- `rule = "ctx.auth != null && ctx.auth.id == item.owner.id"`
- **admin-users**: Allows access to admin users
	- `rule = "ctx.auth != null && ctx.auth.role == 'admin'"`


#### Custom Policies

Custom policies can be defined and implemented by the user. **All custom policies must explicitly define a `scope` field** (`request` or `resource`).

For example:

```hcl
custom "trial-users" {
	scope = "request" # required
	description = "Allows access to trial users"
	# Custom logic to be implemented in user code
}
```


#### Policy Groups

Policies can be grouped for reuse, but **all policies in a group must have the same scope**. Cross-scope grouping is not allowed.

```hcl
# NOTE: All policies in a group must have the same scope. Cross-scope grouping is not allowed.
group "write-access" {
	scope = "request"
	description = "Allows access to admin users"
	policies = ["admin-users", "authenticated"]
}
```

### Applying Policies

Policies can be applied at various levels:

- **Global**: Applies to all services and operations
- **Service**: Applies to a specific service (e.g., `service "user" { ... }`)
- **Operation**: Applies to a specific operation (e.g., `operation "like" { ... }`)

You can also use the `apply` block to attach a policy to specific operations and optionally associate a rate limit:

```hcl
apply "policy" "only-owner" {
	to_operations = ["update", "delete"]
	rate_limit = "user-total"
}
```

### Example Policy Block (from `services.hcl`)

```hcl
policies {
	mode = "deny-by-default"
	precedence = "deny-over-allow"
	short_circuit = true

	policy "authenticated" { ... }
	policy "guest-readonly" { ... }
	policy "only-owner" { ... }
	policy "admin-users" { ... }
	custom "trial-users" { ... }
	group "write-access" { ... }
}
```

---

|----------
Do not touch below
----------|

| Feature Name                          | Status      | Complexity   | Notes                                                                                  |
|---------------------------------------|-------------|--------------|----------------------------------------------------------------------------------------|
| Route Definitions                     | Supported   | Medium       | RESTful, nested, custom operations                                                     |
| Policies (Authorization)              | Supported   | Medium       | Preset, custom, group policies                                                         |
| Rate Limiting                         | Supported   | Medium       | Global, per-service, custom, multiple strategies                                       |
| CORS Configuration                    | Supported   | Low          | Global and per-service                                                                 |
| Middlewares                           | Supported   | Low          | Global and per-service                                                                 |
| CRUD Operations                       | Supported   | Low          | Standard and batch, soft delete                                                        |
| Pagination, Sorting, Filtering        | Supported   | Low          | Defaults, can be customized                                                            |
| Service Exposure/Path Customization   | Supported   | Low          | Expose/hide, custom paths                                                              |
| Health Check Endpoint                 | Supported   | Low          | Simple GET endpoint                                                                    |
| Webhooks & Event Triggers             | Not yet     | High         | Emit events or call webhooks on operations                                             |
| Request/Response Validation           | Not yet     | Medium       | Schema-based, prevents invalid data                                                    |
| Input/Output DTOs                     | Not yet     | Medium       | Custom request/response shapes                                                         |
| Role-Based Access Control (RBAC)      | Not yet     | Medium       | More granular roles/permissions                                                        |
| API Versioning                        | Not yet     | Medium       | Support multiple API versions                                                          |
| OpenAPI/Swagger Generation            | Not yet     | Medium       | Auto-generate docs and SDKs                                                            |
| File Upload/Download Support          | Not yet     | Medium       | Handle file uploads/downloads                                                          |
| Custom Error Handling                 | Not yet     | Medium       | Centralized, customizable error responses                                              |
| Scheduled Jobs/Tasks                  | Not yet     | High         | Background jobs, cron-like features                                                    |
| Multi-tenancy                         | Not yet     | High         | Tenant isolation/scoping                                                               |
| GraphQL Endpoint Generation           | Not yet     | High         | Generate GraphQL endpoints                                                             |
| Service Dependencies/Orchestration    | Not yet     | High         | Define dependencies, orchestrate workflows                                             |
| Rate Limit Quotas per User/Plan       | Not yet     | Medium       | Dynamic rate limits based on user plans                                                |
| Audit Logging                         | Not yet     | Medium       | Automatic change/access logging                                                        |
| API Key Management                    | Not yet     | Medium       | Issue, rotate, validate API keys                                                       |
| Localization/Internationalization     | Not yet     | Medium       | Multi-language support                                                                 |
| Request/Response Transformation Hooks | Not yet     | Medium       | Custom logic before/after handlers                                                     |
| Soft/Hard Delete Toggle per Operation | Not yet     | Low          | More granular delete control                                                           |
| Service/Operation Deprecation Notices | Not yet     | Low          | Mark deprecated services/operations                                                    |
| Integration with Monitoring/Tracing   | Not yet     | Medium       | Metrics, tracing, external monitoring                                                  |