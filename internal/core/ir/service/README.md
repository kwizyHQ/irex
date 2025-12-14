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