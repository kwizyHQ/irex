template "index.ts.tpl" {
  data   = "models:index"
  output = "src/index.ts"
  mode   = "single"
}

template "model.ts.tpl" {
  data   = "model:full"
  output = "src/models/{{ lower .Name }}.ts"
  mode   = "per-item"
}
