package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	"go-learning/practise/cobra-practise/cmd/apply"
	"go-learning/practise/cobra-practise/cmd/plugin"
)

// PluginHandler is capable of parsing command line arguments
// and performing executable filename lookups to search
// for valid plugin files, and execute found plugins.
type PluginHandler interface {
	// exists at the given filename, or a boolean false.
	// Lookup will iterate over a list of given prefixes
	// in order to recognize valid plugin filenames.
	// The first filepath to match a prefix is returned.
	Lookup(filename string) (string, bool)
	// Execute receives an executable's filepath, a slice
	// of arguments, and a slice of environment variables
	// to relay to the executable.
	Execute(executablePath string, cmdArgs, environment []string) error
}
type PixiuOptions struct {
	PluginHandler PluginHandler
	Arguments     []string
	ConfigFlags   *genericclioptions.ConfigFlags

	genericclioptions.IOStreams
}

// DefaultPluginHandler implements PluginHandler
type DefaultPluginHandler struct {
	ValidPrefixes []string
}

func NewDefaultPluginHandler(validPrefixes []string) *DefaultPluginHandler {
	return &DefaultPluginHandler{
		ValidPrefixes: validPrefixes,
	}
}

func (h *DefaultPluginHandler) Lookup(filename string) (string, bool) {
	// TODO
	return "", false
}

func (h *DefaultPluginHandler) Execute(executablePath string, cmdArgs, environment []string) error {
	return nil
}

func NewDefaultPixiuCommand() *cobra.Command {
	return NewDefaultPixiuctlCommandWithArgs(PixiuOptions{
		PluginHandler: NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     os.Args,
		ConfigFlags:   genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag(),
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

// NewDefaultPixiuCommand creates the `pixiuctl` command with default arguments
func NewDefaultPixiuctlCommandWithArgs(o PixiuOptions) *cobra.Command {
	cmd := NewPixiuCommand(o)

	if o.PluginHandler == nil {
		return cmd
	}

	// 预留命令行插件，暂时不做实现
	if len(o.Arguments) > 1 {
	}

	return cmd
}

// NewPixiuCommand 创建 `pixiuctl` 命令行和它的子命令
func NewPixiuCommand(o PixiuOptions) *cobra.Command {
	//warningHandler := rest.NewWarningWriter(o.IOStreams.ErrOut, rest.WarningWriterOptions{Deduplicate: true, Color: term.AllowsColorOutput(o.IOStreams.ErrOut)})
	warningsAsErrors := false

	// Parent command to which all subcommands are added
	cmds := &cobra.Command{
		Use:   "pixiuctl",
		Short: i18n.T("pixiuctl controls the Pixiu cluster manager"),
		Long: templates.LongDesc(`
      pixiuctl controls the Pixiu cluster manager.

      Find more information at:
            https://github.com/caoyingjunz/go-learning`),
		Run: runHelp,
		//TODO： 执行前后钩子
	}

	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	flags := cmds.PersistentFlags()

	// 通过 addFlag 追加
	flags.BoolVar(&warningsAsErrors, "warnings-as-errors", warningsAsErrors, "Treat warnings received from the server as errors and exit with a non-zero exit code")

	pixiuConfigFlags := o.ConfigFlags
	pixiuConfigFlags.AddFlags(flags)
	matchVersionPixiuConfigFlags := cmdutil.NewMatchVersionFlags(pixiuConfigFlags)

	f := cmdutil.NewFactory(matchVersionPixiuConfigFlags)

	cmdGroups := templates.CommandGroups{
		{
			Message:  "Basic Commands (Beginner):",
			Commands: []*cobra.Command{},
		},
		{
			Message: "Deploy Commands:",
			Commands: []*cobra.Command{
				apply.NewCmdApply("pixiuctl", f, o.IOStreams),
			},
		},
	}

	cmdGroups.Add(cmds)

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	fmt.Println("command help")
	//cmd.Help()
}
