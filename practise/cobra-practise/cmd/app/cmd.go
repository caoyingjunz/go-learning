package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	"go-learning/practise/cobra-practise/cmd/apply"
	"go-learning/practise/cobra-practise/cmd/create"
	"go-learning/practise/cobra-practise/cmd/plugin"
	pcmdutil "go-learning/practise/cobra-practise/cmd/util"
)

var (
	defaultKubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
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

// ConfigFlags composes the set of values necessary
// for obtaining pixiu client config
type ConfigFlags struct {
	Kubeconfig *string
	Name       *string
	Namespace  *string

	usePersistentConfig bool
	//TODO
}

// NewConfigFlags returns ConfigFlags with default values set
func NewConfigFlags(usePersistentConfig bool) *ConfigFlags {
	return &ConfigFlags{
		Kubeconfig: stringptr(defaultKubeConfig),
		Name:       stringptr(""),
		Namespace:  stringptr(""),

		usePersistentConfig: usePersistentConfig,
	}
}

func (f *ConfigFlags) WithDefaultNamespaceFlag() *ConfigFlags {
	f.Namespace = stringptr("default")
	return f
}

func stringptr(val string) *string {
	return &val
}

const (
	flagName      = "name"
	flagNamespace = "namespace"
)

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.Kubeconfig != nil {
		flags.StringVar(f.Kubeconfig, "kubeconfig", *f.Kubeconfig, "Path to the kubeconfig file to use for CLI requests.")
	}
	if f.Name != nil {
		flags.StringVar(f.Name, flagName, *f.Name, "Name to impersonate for the operation")
	}
	if f.Namespace != nil {
		flags.StringVar(f.Namespace, flagNamespace, *f.Namespace, "Namespace")
	}

	// TODO: 其他的自定义配置
}

type PixiuOptions struct {
	PluginHandler PluginHandler
	Arguments     []string
	ConfigFlags   *ConfigFlags

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
		ConfigFlags:   NewConfigFlags(true).WithDefaultNamespaceFlag(),
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

// NewDefaultPixiuCommand creates the `pixiuctl` command with default arguments
func NewDefaultPixiuctlCommandWithArgs(o PixiuOptions) *cobra.Command {
	cmd := NewPixiuCommand(o)

	if o.PluginHandler == nil {
		// TODO: 后续考虑加入plugin的实现
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
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if warningsAsErrors {
				fmt.Println("demo warningsAsErrors")
			}
			return nil
		},
	}

	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	// 通过 addFlag 追加
	flags := cmds.PersistentFlags()

	flags.BoolVar(&warningsAsErrors, "warnings-as-errors", warningsAsErrors, "Treat warnings received from the server as errors and exit with a non-zero exit code")

	configFlags := o.ConfigFlags
	configFlags.AddFlags(flags)

	f := pcmdutil.NewFactory(*configFlags.Kubeconfig)

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				create.NewCmdCreate(f, o.IOStreams),
			},
		},
		{
			Message: "Advanced Commands:",
			Commands: []*cobra.Command{
				apply.NewCmdApply(f, o.IOStreams),
			},
		},
	}

	groups.Add(cmds)

	filters := []string{"options"}
	templates.ActsAsRootCommand(cmds, filters, groups...)

	// Stop warning about normalization of flags. That makes it possible to
	// add the klog flags later.
	cmds.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
