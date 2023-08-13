package template

const (
	ServiceTemplate = `---
apiVersion: v1
kind: Service
metadata:
  name: {{ .UserName }}
spec:
{{- if gt (len .Emails) 0 }}
  emails:
{{- range .Emails }}
    - {{ . }}
{{- end }}
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

{{- .Data | nindent 4 }}
`
)
