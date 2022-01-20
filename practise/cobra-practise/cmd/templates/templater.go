package templates

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
	"k8s.io/kubectl/pkg/util/term"
)

type FlagExposer interface {
	ExposeFlags(cmd *cobra.Command, flags ...string) FlagExposer
}

func ActsAsRootCommand(cmd *cobra.Command, filters []string, groups ...templates.CommandGroup) FlagExposer {
	if cmd == nil {
		panic("nil root command")
	}
	templater := &templater{
		RootCmd:       cmd,
		UsageTemplate: MainUsageTemplate(),
		HelpTemplate:  MainHelpTemplate(),
		CommandGroups: groups,
		Filtered:      filters,
	}
	cmd.SetFlagErrorFunc(templater.FlagErrorFunc())
	cmd.SilenceUsage = true
	cmd.SetUsageFunc(templater.UsageFunc())
	cmd.SetHelpFunc(templater.HelpFunc())
	return templater
}

type templater struct {
	UsageTemplate string
	HelpTemplate  string
	RootCmd       *cobra.Command
	templates.CommandGroups
	Filtered []string
}

func (templater *templater) HelpFunc() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		t := template.New("help")
		t.Funcs(templater.templateFuncs())
		template.Must(t.Parse(templater.HelpTemplate))
		out := term.NewResponsiveWriter(c.OutOrStdout())
		err := t.Execute(out, c)
		if err != nil {
			c.Println(err)
		}
	}
}

func (templater *templater) ExposeFlags(cmd *cobra.Command, flags ...string) FlagExposer {
	cmd.SetUsageFunc(templater.UsageFunc(flags...))
	return templater
}

func (templater *templater) templateFuncs(exposedFlags ...string) template.FuncMap {
	fmt.Println("help")
	return template.FuncMap{}
}

func (templater *templater) UsageFunc(exposedFlags ...string) func(*cobra.Command) error {
	return func(c *cobra.Command) error {
		t := template.New("usage")
		t.Funcs(templater.templateFuncs(exposedFlags...))
		template.Must(t.Parse(templater.UsageTemplate))
		out := term.NewResponsiveWriter(c.OutOrStderr())
		return t.Execute(out, c)
	}
}

func (templater *templater) FlagErrorFunc(exposedFlags ...string) func(*cobra.Command, error) error {
	return func(c *cobra.Command, err error) error {
		c.SilenceUsage = true
		switch c.CalledAs() {
		case "options":
			return fmt.Errorf("%s\nRun '%s' without flags.", err, c.CommandPath())
		default:
			return fmt.Errorf("%s\nSee '%s --help' for usage.", err, c.CommandPath())
		}
	}
}

// MainUsageTemplate if the template for 'usage' used by most commands.
func MainUsageTemplate() string {
	return ""
}

func MainHelpTemplate() string {
	return ""
}
