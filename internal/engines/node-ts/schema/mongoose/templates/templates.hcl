template "adaptor.ts.tpl" {
  data   = "schema:index"
  output = "models/dl.types.ts" // with respect to generated folder defined in irex.hcl
  mode   = "per-item"
}

template "index.ts.tpl" {
  data   = "schema:index"
  output = "models/index.ts"
  mode   = "single"
}

template "model.ts.tpl" {
  data   = "schema:model"
  output = "models/{{ lower .Name }}.ts" // with respect to generated folder defined in irex.hcl
  mode   = "per-item"
}
