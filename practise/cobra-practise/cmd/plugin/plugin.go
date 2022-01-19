package plugin

import (
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	pluginListLong = templates.LongDesc(i18n.T(`
		List all available plugin files on a user's PATH.

		Available plugin files are those that are:
		- executable
		- anywhere on the user's PATH
		- begin with "pixiuctl-"
	`))

	ValidPluginFilenamePrefixes = []string{"pixiuctl"}
)
