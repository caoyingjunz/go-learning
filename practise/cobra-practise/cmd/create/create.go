package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	pcmdutil "go-learning/practise/cobra-practise/cmd/util"
)

var (
	createLong = templates.LongDesc(i18n.T(`
		Create a resource from a file or from stdin.

		JSON and YAML formats are accepted.`))

	createExample = templates.Examples(i18n.T(`
		# Create a pod using the data in pod.json
		pixiuctl create -f ./create.json`))
)

type CreateOptions struct {
	Raw              string
	EditBeforeCreate bool

	genericclioptions.IOStreams
}

// ValidateArgs makes sure there is no discrepency in command options
func (o *CreateOptions) ValidateArgs(cmd *cobra.Command, args []string) error {

	return nil
}

// Complete completes all the required options
func (o *CreateOptions) Complete(cmd *cobra.Command) error {
	fmt.Println("test create complete raw:", o.Raw)

	return nil
}

func (o *CreateOptions) RunCreate(cmd *cobra.Command) error {
	fmt.Println("run create edit", o.EditBeforeCreate, "raw", o.Raw)
	return nil
}

func NewCreateOptions(ioStreams genericclioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdCreate(f pcmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "create -f FILENAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Create a resource from a file or from stdin"),
		Long:                  createLong,
		Example:               createExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.ValidateArgs(cmd, args))
			cmdutil.CheckErr(o.RunCreate(cmd))
		},
	}

	// 绑定参数
	cmd.Flags().BoolVar(&o.EditBeforeCreate, "edit", o.EditBeforeCreate, "Edit the API resource before creating")
	cmd.Flags().StringVar(&o.Raw, "raw", o.Raw, "Raw URI to POST to the server.  Uses the transport specified by the kubeconfig file.")

	// create subcommands
	cmd.AddCommand(NewCmdCreateService(ioStreams))
	return cmd
}
