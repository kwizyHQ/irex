# IREX – Copilot Instructions (2025)

These concise guidelines help Copilot work effectively in the IREX project, an AST-driven backend code generator for Node.js (Fastify/Express), Mongoose, Knex, and TypeScript/JavaScript.

---

## Project Structure

- `cmd/` – CLI entrypoints (`irex`, `irex-dev`)
- `internal/` – Core engine: AST, CLI, codegen, templates
- `internal/core/ast/` – Config, schema, service ASTs & templates
- `internal/engines/node-ts/` – TypeScript codegen (bootstrap, schema, service)
- `extensions/vscode/` – VS Code extension
- `examples/` – Example generated apps
- `docs/` – Documentation
- `temp/` – Temporary/generated files

---

## Templates & Codegen

- Template files: `.ts.tpl`, `.js.tpl` (TypeScript/JavaScript + Go template syntax)
- Keep templates minimal, structural, and free of business logic
- Use idiomatic, modern TypeScript/JavaScript
- Use Go template expressions only for dynamic content (e.g. `{{ .Name }}`)

---

## Generated Code

- Fastify: Use `fastify.register()`, async/await, minimal controllers
- Express: Use `Router()`
- Mongoose: Type-safe schemas, timestamps, prefer `.lean()`
- Knex: Use query builder, avoid raw SQL

---

## CLI & AST

- CLI commands: `irex init`, `irex dev`, `irex build`, `irex generate model|service|workflow <Name>`
- AST/spec: Strict, minimal Go structs, JSON-friendly, extendable

---

## Code Style

- TypeScript: ES modules, async/await, strong typing, interfaces, named exports
- Go: Idiomatic, small functions, minimal concurrency (core engine only)

---

## Hooks & Extensibility

- Hooks are user-defined, never overwritten
- Generated files import hooks but do not embed logic

---

## Restrictions

- Do not suggest code for unsupported languages (Rust, C++, PHP, etc.)
- Do not add business logic to templates or generated files

---

## Copilot Goals

- Help with templates, AST→template mapping, CLI, codegen, docs, and extension structure

---

## If Unsure

- Prefer minimal, clean, idiomatic TypeScript and Go
- Keep template and AST/spec structure predictable and simple