package template

const (
	CommonService = `common:
  apiVersion: {{ apiVersion }}
  kind: "{{ kind }}"
  metadata: { metadata }
`
)
