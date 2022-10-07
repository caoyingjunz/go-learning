package template

const (
	ServiceTemplate = `---
apiVersion: v1
kind: Service
metadata:
  name: {{ .UserName }}
spec:
  emails:
{{- range .Emails }}
    - {{ . }}
{{- end }}
  selector:
{{- with .Friends }}
{{- range . }}
    app: {{ .Name }}
{{- end }}
{{- end }}

{{- range $m, $v := .Mods }}
    {{ $m }}: {{ $v }}
{{- end }}
`
)
