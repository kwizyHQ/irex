import mongoose from "mongoose";

export const {{ .Name }}Schema = new mongoose.Schema({
{{- range .Fields }}
  {{ .Name }}: {
    {{- if .IsArray }}
    type: [{{ .MongooseType }}],
    {{- else }}
    type: {{ .MongooseType }},
    {{- end }}

    {{- if .Required }} required: true,{{ end }}
    {{- if .Unique }} unique: true,{{ end }}
    {{- if .Trim }} trim: true,{{ end }}
    {{- if .Min }} min: {{ .Min }},{{ end }}
    {{- if .Max }} max: {{ .Max }},{{ end }}
    {{- if .MinLength }} minlength: {{ .MinLength }},{{ end }}
    {{- if .MaxLength }} maxlength: {{ .MaxLength }},{{ end }}
    {{- if .Match }} match: /{{ .Match }}/,{{ end }}
    {{- if .Default }} default: {{ ctyParse .Default }},{{ end }}
  },
{{- end }}
}, {
  collection: "{{ .Config.Collection }}",
  timestamps: {{ .Config.Timestamps }},
  strict: {{ .Config.Strict }},
  versionKey: {{ .Config.VersionKey }},
  autoIndex: {{ .Config.AutoIndex }},
  autoCreate: {{ .Config.AutoCreate }},
  minimize: {{ .Config.Minimize }},
  strictQuery: {{ .Config.StrictQuery }},
  toJSON: { getters: {{ .Config.ToJSONGetters }} }
});

{{- range .Indexes }}
{{ $.Name }}Schema.index(
  { {{ range $i, $f := .Fields }}{{ if $i }}, {{ end }}{{ $f }}: 1{{ end }} },
  { unique: {{ .Unique }} }
);
{{- end }}

export default mongoose.model("{{ .Name }}", {{ .Name }}Schema);

