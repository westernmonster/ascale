{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}

{{/*
Inject extra environment vars in the format key:value, if populated
*/}}
{{- define "app.exEnv" -}}
{{- if .exEnv -}}
{{- range $key, $value := .exEnv }}
- name: {{ $key }}
  value: {{ $value | quote }}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "app.joinListWithComma" -}}
{{- $local := dict "first" true -}}
{{- range $k, $v := . -}}{{- if not $local.first -}},{{- end -}}"{{- $v -}}"{{- $_ := set $local "first" false -}}{{- end -}}
{{- end -}}




