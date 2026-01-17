# ─────────────────────────────────────────────
# Core Application
# ─────────────────────────────────────────────

template "app.ts.tpl" {
  data   = "service:app"
  output = "app.ts"
  mode   = "single"
}

# ─────────────────────────────────────────────
# Routes
# ─────────────────────────────────────────────

template "route.ts.tpl" {
  data   = "service:routes"
  output = "routes/{{ lower .Name }}.route.ts"
  mode   = "per-item"
}

template "routes.index.ts.tpl" {
  data   = "service:routes_index"
  output = "routes/index.ts"
  mode   = "single"
}

# ─────────────────────────────────────────────
# Controllers
# ─────────────────────────────────────────────

template "controller.ts.tpl" {
  data   = "service:controllers"
  output = "controllers/{{ lower .Name }}.controller.ts"
  mode   = "per-item"
}

template "controllers.index.ts.tpl" {
  data   = "service:controllers_index"
  output = "controllers/index.ts"
  mode   = "single"
}
