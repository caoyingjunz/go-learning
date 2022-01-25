package create

import (
	"fmt"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
)

var (
	serviceClusterIPLong = templates.LongDesc(i18n.T(`
    Create a ClusterIP service with the specified name.`))

	serviceClusterIPExample = templates.Examples(i18n.T(`
    # Create a new ClusterIP service named my-cs (in headless mode)
    pixiuctl create service clusterip my-cs --clusterip="None"`))
)

// NewCmdCreateService is a macro command to create a new service
func NewCmdCreateService(ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewServiceOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   i18n.T("Create a ClusterIP service"),
		Long:    serviceClusterIPLong,
		Example: serviceClusterIPExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.Name, "name", o.Name, "usage name cluster ip")
	cmd.Flags().StringVar(&o.ClusterIP, "clusterip", o.ClusterIP, "usage cluster ip")

	return cmd
}

type ServiceOptions struct {
	Name      string
	Namespace string

	ClusterIP string

	genericclioptions.IOStreams
}

func (o *ServiceOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate validates required fields are set to support structured generation
func (o *ServiceOptions) Validate() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("name must be specified")
	}

	return nil
}

func (o *ServiceOptions) Run() error {
	fmt.Println(fmt.Sprintf("create service: namespace %s, name %s, clusterip: %s", o.Namespace, o.Name, o.ClusterIP))
	return nil
}

func NewServiceOptions(ioStreams genericclioptions.IOStreams) *ServiceOptions {
	return &ServiceOptions{
		IOStreams: ioStreams,
	}
}
