# IREX – Copilot Instructions

These are high-level guidelines for GitHub Copilot when working in the IREX project. The IREX framework is an AST-based backend code generator that currently targets:

- Fastify (primary backend framework)
- Express (secondary)
- Mongoose (MongoDB)
- Knex (SQL databases)
- TypeScript (preferred)
- JavaScript (allowed)

Future runtimes like Go, Rust, C++, PHP, etc. will come later, but Copilot should not suggest them yet.

---

## 1. Project architecture

IREX follows an AST (spec/DSL) → Template → Output pipeline.

Key folders:
- `internal/*` → AST (spec/DSL) builder, parser, generator, renderer, cli etc.
- `cmd/irex/` → Build CLI commands - these commands will be used by end users
- `cmd/irex-dev/` → Development commands - used during IREX generator development
- `docs/` → Documentation
- `gen/` → Generated example applications

Copilot should help maintain consistent template structure, clean AST→render logic, modular generators (Fastify/Express, Mongoose/Knex), and a TypeScript-first approach.

---

## 2. Template file rules

Template files use extensions like `.ts.tpl` and `.js.tpl`. They must contain valid TypeScript/JavaScript with Go-template expressions when needed.

Guidelines:

- Generate idiomatic, minimal, clean TS/JS.
- Never include business logic inside templates.
- Keep templates pure and focused on structure.
- Use Go template syntax only when necessary, for example:

  {{ .Name }}
  {{ range .Fields }}
  {{ end }}

---

## 3. Generated code rules

Generated code should be simple, idiomatic and framework-appropriate.

Fastify

- Use `fastify.register()` for modules.
- Route files should look like:

```ts
export default async function (fastify) {
  fastify.get(...)
  fastify.post(...)
}
```

Prefer async/await and keep controllers minimal. For Express use `Router()`.

Mongoose

- Always define schemas with type-safe fields and timestamps.
- Use `.lean()` for queries when possible.

Knex

- Use query builders and avoid raw SQL unless necessary.

---

## 4. Hook system guidelines

Hooks must be user-defined and never overwritten. Provide clean extension points and avoid embedding logic in generated files. Respect import patterns for hook files and only suggest stubs when asked.

Example:

```ts
import * as UserHooks from "../../hooks/user.hooks";
```

---

## 5. IREX CLI conventions

Common CLI commands:

- `irex init`
- `irex dev`
- `irex build`
- `irex generate model <ModelName>`
- `irex generate service <ServiceName>`
- `irex generate workflow <WorkflowName>`

Copilot should follow this pattern, keep CLI commands modular, and use Cobra-like structures in Go when appropriate.

---

## 6. AST/spec guidelines

AST/spec is represented as structured Go data. Keep definitions strict, typed, minimalistic, JSON-friendly and extendable (workflows, models, services).

Example:

```go
type ASTModel struct {
  Name   string
  Fields []ASTField
}
```

---

## 7. Code style guidelines

TypeScript:

- Use ES modules (import … from).
- Prefer async/await.
- Strong typing; avoid `any`.
- Use interfaces for DTOs and prefer named exports.

Go (core engine only — `internal/*`):

- Keep Go code idiomatic.
- Avoid unnecessary concurrency.
- Write small, simple functions.

---

## 8. Restrictions / avoid

Do not suggest application layer code in languages outside the current scope (Rust, C++, PHP, Python, etc.).

Do not mix template syntax incorrectly, add business logic inside generated files, or suggest unnecessary abstractions.

---

## 9. Goals for Copilot assistance

Help with:

- Writing clean template files.
- Mapping AST/spec → template placeholders.
- Implementing CLI commands.
- Structuring Fastify/Express modules.
- Writing clean TypeScript/Javascript services.
- Writing AST/spec builders/parsers in Go.
- Documentation (`docs/*.md`).

---

## 10. If unsure

Prefer simplicity, minimalism, clean TypeScript, predictable AST/spec formats and consistent template structures.

---

If you want, I can also generate additional repo helpers such as a Copilot contextual prompt, a Copilot Chat persona, a `CONTRIBUTING.md`, or a template linting ruleset — tell me which and I can add one.

## 11. Run go commands

### irex commands

IREX user commands: These are used by end users of the IREX framework.
```
irex init # Initialize a new IREX project
irex dev # Start development server with hot-reloading
irex build # Build the IREX project
```

### irexd commands

IREX development commands: These are used when developing IREX itself.

```
irexd init # Initialize IREX development environment
irexd watch # Watch for changes (In project + go files) and rebuild IREX 