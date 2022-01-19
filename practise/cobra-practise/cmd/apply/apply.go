package apply

import (
	"fmt"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	applyLong = templates.LongDesc(i18n.T(`
		Apply a configuration to a resource by file name or stdin.
		The resource name must be specified. This resource will be created if it doesn't exist yet.
		To use 'apply', always create the resource initially with either 'apply' or 'create --save-config'.

		JSON and YAML formats are accepted.

		Alpha Disclaimer: the --prune functionality is not yet complete. Do not use unless you are aware of what the current state is. See https://issues.k8s.io/34274.`))

	applyExample = templates.Examples(i18n.T(`
		# Apply the configuration in pod.json to a pod
		pixiuctl apply -f ./test.json`))
)

// NewCmdApply create the `apply` command
func NewCmdApply(baseName string) *cobra.Command {
	flags := NewApplyFlags()

	cmd := &cobra.Command{
		Use:                   "apply (-f FILENAME | -k DIRECTORY)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Apply a configuration to a resource by file name or stdin"),
		Long:                  applyLong,
		Example:               applyExample,
		Run: func(cmd *cobra.Command, args []string) {
			o, err := flags.ToOptions(cmd, baseName, args)
			cmdutil.CheckErr(err)
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

// ApplyFlags directly reflect the information that CLI is gathering via flags.  They will be converted to Options, which
// reflect the runtime requirements for the command.  This structure reduces the transformation to wiring and makes
// the logic itself easy to unit test
type ApplyFlags struct {
	genericclioptions.IOStreams
}

// NewApplyFlags returns a default ApplyFlags
func NewApplyFlags() *ApplyFlags {
	return &ApplyFlags{}
}

// ApplyOptions defines flags and other configuration parameters for the `apply` command
type ApplyOptions struct {
	Kubeconfig string

	Name      string
	Namespace string

	genericclioptions.IOStreams
}

// ToOptions converts from CLI inputs to runtime inputs
func (flags *ApplyFlags) ToOptions(cmd *cobra.Command, baseName string, args []string) (*ApplyOptions, error) {
	kubeConfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return nil, err
	}
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return nil, err
	}
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}

	o := &ApplyOptions{
		Kubeconfig: kubeConfig,
		Namespace:  namespace,
		Name:       name,

		IOStreams: flags.IOStreams,
	}

	return o, nil
}

// Run executes the `apply` command.
func (o *ApplyOptions) Run() error {
	// TODO: run with options
	fmt.Println(fmt.Sprintf("run apply command with Kubeconfig: %s, Namespace: %s, Name: %s", o.Kubeconfig, o.Namespace, o.Name))
	return nil
}
