# IREX — Copilot Instructions

These concise instructions help AI coding agents contribute effectively to IREX (Intermediate Representation Executor). Focus on concrete, discoverable patterns and files so you can be productive immediately.

1. Big picture
- **Purpose:** IREX parses HCL into an IR and generates backend code (Node.js/TypeScript, Fastify/Express, Mongoose/Knex). See root README for overview.
- **Pipeline:** Input HCL → AST (`internal/core/ast`) → semantic/validation (`internal/core/semantic`, `internal/core/validate`) → normalization (`internal/core/normalize`) → IR (`internal/ir`) → engines (`internal/engines/*`) for codegen.

2. Key directories to read first
- **`internal/`**: core engine, AST, pipeline, and transform logic (start here to understand IR shapes). Example: `internal/core/ast/entry.go` and `internal/ir/service.go`.
- **`internal/engines/node-ts/`**: TypeScript codegen and bootstrapping (templates and scaffold). Look at `bootstrap/scaffold.go` and `bootstrap/templates/scaffold/*` for generated app layout.
- **`extensions/vscode/irex/`**: VS Code extension integration and commands — useful for editor UX patterns and tests.
- **`cmd/`**: CLI entrypoints and main programs (`irex`, `irex-dev`) — check `cmd/irex/main.go` and `cmd/irex-dev/main.go` for flags and runtime behavior.

3. Templates and generation patterns
- Templates use Go templates embedded in the engine. Look for `templates/` under engines (for example `internal/engines/node-ts/bootstrap/templates/scaffold/app.ts`).
- Template file extensions: `.ts.tpl`, `.js.tpl`. VS Code mappings exist in `.vscode/settings.json`.
- Keep templates minimal: logic lives in Go code that prepares template data. Inspect `internal/engines/*/*/*.go` for the struct fields passed into templates.

4. Build / dev workflows (how to run & test)
- Build the Go CLI: `go build ./...` at repo root. CLI binaries are produced in place (examples: `irex`, `irex-dev`).
- Run the dev CLI: `go run ./cmd/irex` or `go run ./cmd/irex-dev` with flags used in `main.go` (see `cmd/` files for flags). Use `-h` to list commands.
- Generated app scaffold: `internal/engines/node-ts/bootstrap/scaffold.go` creates `src/app.ts`, `src/vendor/server.ts`, `.env.example`, and `README.md` in generated projects.

5. Project-specific conventions
- Single-source HCL specs drive everything; prefer changes in parsers/ASTs over changing templates when adding features.
- Validation and semantic checks are centralized: `internal/core/validate` and `internal/core/semantic` — add checks there so all engines benefit.
- IR is favored for inter-engine communication: edit `internal/ir` types carefully; they are canonical.

6. Tests and diagnostics
- Diagnostics and validation helpers live in `internal/diagnostics` and are used by validators and pipeline stages.
- Extension tests live under `extensions/vscode/irex/test/` — run via the extension test runner in that folder.

7. Integration points & external deps
- Engines target external runtimes (Node.js, TypeScript). Generated projects may expect Node tooling (`nodemon`, `ts-node`) — see `templates/scaffold/nodemon.json`.
- The VS Code extension uses the generated templates and example workspaces under `extensions/vscode/irex/`.

8. Useful code references (examples to inspect)
- `internal/core/ast/entry.go` — AST entry points and structs.
- `internal/core/semantic/service.go` — semantic checks for services.
- `internal/ir/service.go` — canonical IR service model.
- `internal/engines/node-ts/bootstrap/scaffold.go` — how a scaffold is assembled and which templates are used.
- `extensions/vscode/irex/src/commands/index.js` — extension commands and UX hooks.

9. What to avoid / limits
- Do not add business logic into templates; templates receive prepared structs and should remain presentation-only.
- Avoid editing generated output in engines; instead update the templates or the data model feeding them.

10. Examples of focused tasks you can do
- Add a semantic validation: modify `internal/core/semantic/*` and update tests.
- Add a new field to the IR: update `internal/ir/*` and all engine mappings.
- Improve scaffold templates: modify `internal/engines/node-ts/bootstrap/templates/scaffold/*` and the corresponding `scaffold.go` mappings.

If anything here is unclear or you want more detail about a specific area (build commands, test runners, or engine internals), tell me which part and I'll expand or adjust these instructions.
