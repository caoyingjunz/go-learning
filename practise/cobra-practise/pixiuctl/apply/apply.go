package apply

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	pcmdutil "go-learning/practise/cobra-practise/pixiuctl/util"
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
func NewCmdApply(f pcmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewApplyOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "apply (-f FILENAME | -k DIRECTORY)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Apply a configuration to a pixiu resource by file name or stdin"),
		Long:                  applyLong,
		Example:               applyExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.Run())
		},
	}

	// 绑定
	o.AddFlags(cmd)

	return cmd
}

// ApplyOptions defines flags and other configuration parameters for the `apply` command
type ApplyOptions struct {
	Kubeconfig string

	Name      string
	Namespace string

	genericclioptions.IOStreams
}

func NewApplyOptions(ioStreams genericclioptions.IOStreams) *ApplyOptions {
	return &ApplyOptions{
		IOStreams: ioStreams,
	}
}

func (o *ApplyOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to the kubeconfig file to use for CLI requests.")
	cmd.Flags().StringVar(&o.Name, "name", o.Name, "Name to impersonate for the operation")
	cmd.Flags().StringVar(&o.Namespace, "namespace", o.Namespace, "Namespace to impersonate for the operation")
}

func (o *ApplyOptions) Complete(cmd *cobra.Command, args []string) error {
	fmt.Println("test apply complete name:", o.Name)

	return nil
}

// Just a demo
func (o *ApplyOptions) Validate(cmd *cobra.Command) error {
	//if o.Namespace == "" {
	//	return fmt.Errorf("invalied namespace")
	//}

	return nil
}

// Run executes the `apply` command.
func (o *ApplyOptions) Run() error {
	// TODO: run with options
	fmt.Println(fmt.Sprintf("run apply command with Kubeconfig: %s, Namespace: %s, Name: %s", o.Kubeconfig, o.Namespace, o.Name))
	return nil
}
